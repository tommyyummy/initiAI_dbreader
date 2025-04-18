package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func httpConfigProxyHandler(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	action := r.URL.Query().Get("action")

	address := "http://facesong-alpha.wakeai.tech/api/v1/facesong/internal/ccc"
	if env == "prod" {
		address = "https://facesong.wakeai.tech/api/v1/facesong/internal/ccc"
	}
	if *Mode == "local" {
		address = "http://127.0.0.1:8888/api/v1/facesong/internal/ccc"
	}

	var resp *http.Response
	switch action {
	case "rollback":
		resp, _ = http.Post(fmt.Sprintf("%s/rollback/", address)+r.URL.Query().Get("index"), "application/json", r.Body)
	case "delete", "update":
		resp, _ = http.Post(fmt.Sprintf("%s/update", address), "application/json", r.Body)
	default:
		return
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

type ConfigAll struct {
	Index int            `json:"index"`
	Data  map[string]any `json:"data"`
}

func httpConfigChooseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<title>[Config Center]</title>
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
	<div><a class=redis href='/home/config/view?env=test'>test</a>
	</div>`)

	fmt.Fprintf(w, `
	<div class=cell>
	<div><a class=redis href='/home/config/view?env=prod'>prod</a>
	</div>`)
}

func httpConfigHandler(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	address := "http://facesong-alpha.wakeai.tech/api/v1/facesong/internal/ccc"
	if env == "prod" {
		address = "https://facesong.wakeai.tech/api/v1/facesong/internal/ccc"
	}
	if *Mode == "local" {
		address = "http://127.0.0.1:8888/api/v1/facesong/internal/ccc"
	}
	index := r.URL.Query().Get("index")

	resp, _ := http.Get(fmt.Sprintf("%s/all", address))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	configAll := ConfigAll{}
	json.Unmarshal(body, &configAll)
	curIndex := configAll.Index

	if index != "" {
		_, err := strconv.Atoi(index)
		if err == nil {
			resp, _ = http.Get(fmt.Sprintf("%s/index/%s", address, index))
			defer resp.Body.Close()
			body, _ = ioutil.ReadAll(resp.Body)
			//fmt.Println(string(body))

			json.Unmarshal(body, &configAll)
		}
	}

	data := struct {
		CurIndex int
		Env      string
		JsonStr  string
		Host     string
	}{curIndex, env, string(body), *Host}
	templ, err := template.ParseFiles("./config/index.html")

	err = templ.Execute(w, data)
	if err != nil {
		fmt.Println("Error = ", err)
	}
}
