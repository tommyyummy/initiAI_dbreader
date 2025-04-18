package main

import (
	"flag"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
)

var (
	Mode = flag.String("mode", "formal", "")
	Host = flag.String("host", "127.0.0.1:3000", "js http call host")
)

func main() {
	flag.Parse()
	fmt.Printf("Mode: %s\n", *Mode)
	fmt.Printf("Host: %s\n", *Host)

	initMongoClient()
	initRedisClient()

	http.Handle("/home", gziphandler.GzipHandler(http.HandlerFunc(httpHomeHandler)))

	http.Handle("/home/redis", gziphandler.GzipHandler(http.HandlerFunc(httpRedisHandler)))
	http.Handle("/home/redis/list", gziphandler.GzipHandler(http.HandlerFunc(httpListHandler)))
	http.Handle("/home/redis/search", gziphandler.GzipHandler(http.HandlerFunc(httpSearchHandler)))
	http.Handle("/home/redis/detail", gziphandler.GzipHandler(http.HandlerFunc(httpDetailHandler)))

	http.Handle("/home/mongo", gziphandler.GzipHandler(http.HandlerFunc(httpMongoHandler)))
	http.Handle("/home/mongo/collections", gziphandler.GzipHandler(http.HandlerFunc(httpCollectionsHandler)))
	http.Handle("/home/mongo/indexes", gziphandler.GzipHandler(http.HandlerFunc(httpIndexesHandler)))

	http.Handle("/home/config", gziphandler.GzipHandler(http.HandlerFunc(httpConfigChooseHandler)))
	http.Handle("/home/config/view", gziphandler.GzipHandler(http.HandlerFunc(httpConfigHandler)))
	http.Handle("/config/proxy", gziphandler.GzipHandler(http.HandlerFunc(httpConfigProxyHandler)))

	http.Handle("/home/config_v2", gziphandler.GzipHandler(http.HandlerFunc(httpConfigChooseV2Handler)))
	http.Handle("/home/config_v2/view", gziphandler.GzipHandler(http.HandlerFunc(httpConfigV2Handler)))
	http.Handle("/config_v2/proxy", gziphandler.GzipHandler(http.HandlerFunc(httpConfigV2ProxyHandler)))

	fs := http.FileServer(http.Dir("config"))
	http.Handle("/config/", http.StripPrefix("/config/", fs))

	port := ":3000"
	fmt.Println("Server is running on port" + port)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(port, nil))
}
