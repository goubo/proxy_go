package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var port int

func ping(w http.ResponseWriter, req *http.Request) {
	//如果有ip,检查指定ip是否通
	//没有ip ,直接返回pong
	ip := req.URL.Query().Get("ip")
	if ip == "" {
		_, _ = w.Write([]byte("pong"))
	} else {
		port2 := req.URL.Query().Get("port")
		if port2 == "" {
			port2 = strconv.Itoa(port)
		}
		client := http.Client{Timeout: 25 * time.Second}
		resp, err := client.Get(fmt.Sprintf("http://%s:%s%s", ip, port2, req.URL.Path))
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("连接失败"))
			return
		}
		if resp.StatusCode == 200 {
			_, _ = w.Write([]byte("pong"))
		} else {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("连接失败"))
		}
	}
}

func Route(p int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	err := http.ListenAndServe(fmt.Sprintf(":%d", p), mux)
	if err != nil {
		log.Fatalf("http service 启动失败,%v", err)
	}
}
