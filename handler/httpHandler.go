package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var port int

func ping(w http.ResponseWriter, req *http.Request) {
	//如果有ip,检查指定ip是否通
	//没有ip ,直接返回pong
	ip := req.URL.Query().Get("ip")
	client := http.Client{Timeout: 25 * time.Second}
	if ip == "" {
		_, _ = w.Write([]byte("pong"))
	} else {
		resp, err := client.Get(fmt.Sprintf("http://%s:%d%s", ip, port, req.URL.Path))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("连接失败"))
			return
		}
		if resp.StatusCode == 200 {
			w.Write([]byte("pong"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte("连接失败"))
		}
	}
}

func Route(p int) {
	port = p
	http.HandleFunc("/ping", ping)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         fmt.Sprintf(":%d", port),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("http service 启动失败,%v", err)
		os.Exit(-1)
	}

}
