package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
			port2 = string(port)
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
	port = p
	http.HandleFunc("/ping", ping)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("http service 启动失败,%v", err)
		os.Exit(-1)
	}

}

func Soap11(url string, body string) (data string, err error) {
	res, err := http.Post(url, "text/soap; charset=UTF-8", strings.NewReader(body))
	if nil != err {
		fmt.Println("http post err:", err)
		return
	}
	defer res.Body.Close()
	if http.StatusOK != res.StatusCode {
		fmt.Println("WebService soap1.1 request fail, status: %s\n", res.StatusCode)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	return string(result), err

}
