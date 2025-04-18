package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// redis
var (
	rdb = redis.NewClient(&redis.Options{
		//Username: "your redis username",
		Addr: "172.27.177.11:6379",
		//Password: "your redis password",
		PoolSize: runtime.NumCPU() * 20,
		DB:       1,
	})
)

var rdbClient map[string]*redis.Client = make(map[string]*redis.Client)

func initRedisClient() {
	rdbClient["allstar_dev"] = rdb
}

var mdbClient map[string]*mongo.Database = make(map[string]*mongo.Database)

func initMongoClient() {
	mdbTest, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:kh2mMgHfguncXcpZmrrgr1TJpc6272@172.27.177.6:30883/?directConnection=true"))
	mdbClient["allstar_dev"] = mdbTest.Database("allstar_dev")
}

type XMessage struct {
	ID     string `json:"id"`
	Values string `json:"values"`
}

type KVDetail struct {
	Key         string            `json:"key"`
	KeyType     string            `json:"key_type"`
	ValueS      string            `json:"value_s"`
	ValueSL     []string          `json:"value_sl"`
	ValueSD     map[string]string `json:"value_sd"`
	ValueSLL    [][]string        `json:"value_sll"`
	ValueStream []XMessage        `json:"value_stream"`
	TTL         time.Duration     `json:"ttl"`
}

func httpHomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<title>[DB Reader]</title>
	<div style='text-align:left'>
	<h1>DB Reader</h1>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`)

	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=key href='/home/mongo'>mongo</a>
	</div>
	`)
	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=key href='/home/redis'>redis</a>
	</div>
	`)
	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=key href='/home/config'>config center</a>
	</div>
	`)
	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=key href='/home/config_v2'>config center v2</a>
	</div>
	`)
}

func httpRedisHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<title>[RedisReader]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div style='text-align:left'>
	<h1>Redis List</h1>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`)

	for redisName := range rdbClient {
		fmt.Fprintf(w, `
		<div class=cell>
		<div><a class=redis href='/home/redis/list?redis=%s'>%s</a>
		</div>
		`, redisName, redisName)
	}
}

func httpMongoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<title>[MongoReader]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div style='text-align:left'>
	<h1>Mongo List</h1>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`)

	for mongoName := range mdbClient {
		fmt.Fprintf(w, `
		<div class=cell>
		<div><a class=redis href='/home/mongo/collections?mongo=%s'>%s</a>
		</div>
		`, mongoName, mongoName)
	}
}

func httpDetailHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	redisName := r.URL.Query().Get("redis")
	rdb := rdbClient[redisName]

	if key == "" || rdb == nil {
		httpListHandler(w, r)
		return
	}

	keyType := rdb.Type(context.Background(), key).Val()
	ttl := rdb.TTL(context.Background(), key).Val()
	res := KVDetail{Key: key, KeyType: keyType, TTL: time.Duration(ttl.Seconds())}
	switch keyType {
	case "string":
		rawValue := rdb.Get(context.Background(), key).Val()
		res.ValueS = rawValue
	case "set":
		rawValue := rdb.SMembers(context.Background(), key).Val()
		res.ValueSL = rawValue
	case "zset":
		rawValue := rdb.ZRangeWithScores(context.Background(), key, 0, -1).Val()
		temp := [][]string{}
		for _, el := range rawValue {
			temp = append(temp, []string{el.Member.(string), fmt.Sprintf("%f", el.Score)})
		}
		res.ValueSLL = temp
	case "hash":
		rawValue := rdb.HGetAll(context.Background(), key).Val()
		res.ValueSD = rawValue
	case "list":
		rawValue := rdb.LRange(context.Background(), key, 0, -1).Val()
		res.ValueSL = rawValue
	case "stream":
		rawValue := rdb.XRange(context.Background(), key, "-", "+").Val()
		values := make([]XMessage, 0)
		for _, el := range rawValue {
			bytes, _ := json.Marshal(el.Values)
			values = append(values, XMessage{el.ID, string(bytes)})
		}
		res.ValueStream = values
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	resByte, _ := json.Marshal(res)
	w.Write(resByte)
}

func httpCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	mongoName := r.URL.Query().Get("mongo")

	fmt.Fprintf(w, `
	<html>
	<title>[MongoReader]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div class=cell>
	<div><a class=home href='/home/mongo'>mongo list</a>
	</div>
	<div style='text-align:left'>
	<h1>Mongo Reader</h1>
	<h2>%s</h2>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`, mongoName)

	mdb := mdbClient[mongoName]
	if mdb == nil {
		return
	}

	allCollections, err := mdb.ListCollectionNames(context.Background(), bson.M{})
	// fmt.Println(allCollections)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, collectionName := range allCollections {
		fmt.Fprintf(w, `
		<div class=cell>
		<a class=key href='/home/mongo/indexes?mongo=%s&collection=%s'>%s</a>
		</div>`, mongoName, collectionName, collectionName)
	}
}

func httpIndexesHandler(w http.ResponseWriter, r *http.Request) {
	mongoName := r.URL.Query().Get("mongo")
	collection := r.URL.Query().Get("collection")
	field := r.URL.Query().Get("field")
	value := r.URL.Query().Get("value")
	fmt.Println(field)
	fmt.Println(value)

	mdb := mdbClient[mongoName]
	if mdb == nil {
		httpCollectionsHandler(w, r)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, `
	<html>
	<title>[MongoReader]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div class=cell>
	<div><a class=home href='/home/mongo'>mongo list</a>
	</div>
	<div><a class=home href='/home/mongo/collections?mongo=%s'>collections</a>
	</div>
	<div style='text-align:left'>
	<h1>Mongo Reader</h1>
	<h2>%s</h2>
	<h3>%s</h3>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`, mongoName, mongoName, collection)

	fmt.Fprintf(w, `
	<form style='display:inline-block'>
	<div style='text-align:left'>
	<input type=hidden style='width:160px;text-align:center' name=mongo value=%s>
	<input type=hidden style='width:160px;text-align:center' name=collection value=%s>	
	<div>Field</div>
	<input style='width:200px;text-align:left' name=field value=%s>
	<div>Value</div>
	<input style='width:200px;text-align:left' name=value value=%s>
	<div>(only string field and _id supported for now)</div>
	<input type=submit href=/home/mongo/indexes value=search>
	</div>
	`, mongoName, collection, field, value)

	fmt.Fprintf(w, `
	<div class=cell>-------</div>
	`)

	if field == "" || value == "" {
		indexView := mdb.Collection(collection).Indexes()
		cursor, _ := indexView.List(context.Background())
		// if err != nil {
		// 	fmt.Fprintf(w, `
		// 	<div style='text-align:left'>
		// 	<h2>Invalid Collection Name: %s</h2>
		// 	</div>
		// 	`, collection)
		// 	return
		// }
		var result []bson.M
		cursor.All(context.Background(), &result)
		for _, el := range result {
			for k, v := range el {
				fmt.Fprintf(w, `
				<div class=cell>%v: %v</div>
				`, k, v)
				// fmt.Printf("%v: %v\n", k, v)
			}
			fmt.Fprintf(w, `
			<div class=cell>-------</div>
			`)
		}
	} else {
		filter := bson.M{
			field: value,
		}
		if tryUUID, err := uuid.Parse(value); err == nil {
			uuidByte, _ := tryUUID.MarshalBinary()
			filter = bson.M{"$or": []bson.M{bson.M{field: primitive.Binary{byte(0x04), uuidByte}}, bson.M{field: value}}}
		}
		fmt.Println(filter)
		//if field == "_id" {
		//	realValue, err := primitive.ObjectIDFromHex(value)
		//	if err == nil {
		//		filter = bson.M{
		//			field: realValue,
		//		}
		//	}
		//}
		// if err == nil {
		doc := mdb.Collection(collection).FindOne(context.TODO(), filter)
		var result bson.M
		err := doc.Decode(&result)
		fmt.Println(result)
		// fmt.Println(err.Error())
		if err == nil {
			for k, v := range result {
				if b, ok := v.(primitive.Binary); ok {
					realV, err := uuid.FromBytes(b.Data)
					if err == nil {
						fmt.Fprintf(w, `
							<div class=cell>%v: %v</div>
							`, k, realV)
						continue
					}
				}
				if b, ok := v.(primitive.DateTime); ok {
					fmt.Fprintf(w, `
							<div class=cell>%v: %v</div>
							`, k, b.Time())
					continue
				}
				//if b, ok := v.(primitive.M); ok {
				//
				//}
				//fmt.Println(reflect.TypeOf(v))
				//fmt.Printf("%v", v)
				//fmt.Println("")

				fmt.Fprintf(w, `
					<div class=cell>%v: %v</div>
					`, k, v)
				// fmt.Printf("%v: %v\n", k, v)
			}
		}
		// }
	}

}

func doSearchAndList(w http.ResponseWriter, redisName, rawPrefix string, cursor uint64) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	prefix := rawPrefix + "*"

	fmt.Println("input cursor: ", cursor)
	fmt.Println("input redis name: ", redisName)
	fmt.Println("input raw prefix: ", rawPrefix)

	fmt.Fprintf(w, `
	<html>
	<title>[RedisReader]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div class=cell>
	<div><a class=home href='/home/redis'>redis list</a>
	</div>
	<div style='text-align:left'>
	<h1>Redis Reader</h1>
	<h2>%s</h2>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`, redisName)

	rdb := rdbClient[redisName]
	if rdb == nil {
		return
	}

	fmt.Fprintf(w, `
	<form style='display:inline-block'>
	<div style='text-align:left'>
	<input type=hidden style='width:160px;text-align:center' name=redis value=%s>
	<input style='width:160px;text-align:left' name=prefix value=%s><input type=submit formaction=/home/redis/search value=search>
	</div>
	`, redisName, rawPrefix)

	allKeys := []string{}
	size := 20

	start := time.Now()
	for len(allKeys) < 20 {
		keys, next := rdb.Scan(context.Background(), cursor, prefix, int64(size)).Val()
		for _, key := range keys {
			fmt.Fprintf(w, `
			<div class=cell>
			<a class=key href='/home/redis/detail?redis=%s&key=%s'>%s</a>
			</div>`, redisName, key, key)
			allKeys = append(allKeys, key)
		}
		cursor = next
		size *= 2
		if next == 0 {
			break
		}
		maxTime := 1 * time.Second
		if redisName == "test" {
			maxTime = 5 * time.Second
		}
		if time.Since(start) > maxTime {
			fmt.Println("time")
			break
		}
	}

	fmt.Println("output cursor: ", cursor)
	fmt.Println(cursor == 0)

	if cursor == 0 {
		return
	}
	fmt.Fprintf(w, `
	<div style='text-align:left'>
	<input type=hidden style='width:160px;text-align:center' name=cursor value=%s><input type=submit formaction=/home/redis/list name=cursor value=next>
	</div>
	</form>
	`, fmt.Sprintf("%d", cursor))
}

func httpSearchHandler(w http.ResponseWriter, r *http.Request) {
	rawPrefix := r.URL.Query().Get("prefix")
	var cursor uint64 = 0
	redisName := r.URL.Query().Get("redis")

	doSearchAndList(w, redisName, rawPrefix, cursor)
}

func httpListHandler(w http.ResponseWriter, r *http.Request) {
	rawPrefix := r.URL.Query().Get("prefix")

	rawCursor := r.URL.Query().Get("cursor")
	cursor, _ := strconv.ParseUint(rawCursor, 10, 64)

	redisName := r.URL.Query().Get("redis")

	doSearchAndList(w, redisName, rawPrefix, cursor)
}
