package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

func httpConfigV2ProxyHandler(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	action := r.URL.Query().Get("action")

	address := "http://facesong-alpha.wakeai.tech/api/v1/facesong/internal/ccc_v2"
	if env == "prod" {
		address = "https://facesong.wakeai.tech/api/v1/facesong/internal/ccc_v2"
	}
	if *Mode == "local" {
		address = "http://127.0.0.1:8888/api/v1/facesong/internal/ccc_v2"
	}

	var resp *http.Response
	switch action {
	case "namespace_new":
		namespace := r.URL.Query().Get("namespace")
		resp, _ = http.Get(address + "/namespace/" + namespace)
	case "namespace_index":
		namespace := r.URL.Query().Get("namespace")
		index := r.URL.Query().Get("index")
		resp, _ = http.Get(address + "/index/" + namespace + "/" + index)
	case "rollback":
		namespace := r.URL.Query().Get("namespace")
		index := r.URL.Query().Get("index")
		resp, _ = http.Post(fmt.Sprintf("%s/rollback/", address)+namespace+"/"+index, "application/json", r.Body)
	case "update", "delete":
		resp, _ = http.Post(fmt.Sprintf("%s/update", address), "application/json", r.Body)
	default:
		return
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

type ConfigAllV2 struct {
	Index     int            `json:"index"`
	IndexData map[string]any `json:"index_data"`
	//Data       map[string]any `json:"data"`
	Id         string `json:"id"`
	Namespace  string `json:"namespace"`
	Comment    string `json:"comment"`
	CreateTime string `json:"create_time"`
}

func httpConfigChooseV2Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<title>[Config Center V2]</title>
	<div class=cell>
	<div><a class=home href='/home'>home</a>
	</div>
	<div style='text-align:left'>
	<h1>Choose Environments</h1>
	<style>
	.cell { display: inline-block; width: 100em; }
	</style>
	</div>`)

	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=redis href='/home/config_v2/view?env=test'>test</a>
	</div>`)

	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=redis href='/home/config_v2/view?env=prod'>prod</a>
	</div>`)
}

func httpConfigV2Handler(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	address := "http://facesong-alpha.wakeai.tech/api/v1/facesong/internal/ccc_v2"
	if env == "prod" {
		address = "https://facesong.wakeai.tech/api/v1/facesong/internal/ccc_v2"
	}
	if *Mode == "local" {
		address = "http://127.0.0.1:8888/api/v1/facesong/internal/ccc_v2"
	}
	//index := r.URL.Query().Get("index")

	resp, _ := http.Get(fmt.Sprintf("%s/all", address))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	configAll := ConfigAllV2{}
	json.Unmarshal(body, &configAll)
	curIndex := configAll.Index

	//if index != "" {
	//	_, err := strconv.Atoi(index)
	//	if err == nil {
	//		resp, _ = http.Get(fmt.Sprintf("%s/index/%s", address, index))
	//		defer resp.Body.Close()
	//		body, _ = ioutil.ReadAll(resp.Body)
	//		//fmt.Println(string(body))
	//
	//		json.Unmarshal(body, &configAll)
	//	}
	//}

	showJsonStr, _ := json.Marshal(configAll)

	data := struct {
		OverallIndex int
		Env          string
		ShowJsonStr  string
		RealJsonStr  string
		Host         string
	}{curIndex, env, string(showJsonStr), string((body)), *Host}
	templ, err := template.ParseFiles("./config/index_v2.html")

	err = templ.Execute(w, data)
	if err != nil {
		fmt.Println("Error = ", err)
	}
}
