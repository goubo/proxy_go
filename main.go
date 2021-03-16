package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path"
	"path/filepath"
	"proxy/handler"
	"sync"
	"time"
)

var t = flag.String("t", "yaml", "配置文件格式, 支持 json|yaml, 自动读取文件后缀,无后缀需要手动指定")
var c = flag.String("c", "./demo_proxy.yaml", "指定配置文件")
var g = flag.Bool("g", false, "在当前目录生成示例配置文件")
var p = flag.Int("p", 11104, "代理本身服务端口,api处理")

func main() {
	flag.Parse()
	v := viper.New()
	if *g {
		v.SetConfigFile("./demo_proxy." + *t)
		v.Set("ProxyConfig", []handler.ProxyConfig{{
			LocalPort:  43306,
			RemoteIp:   "192.168.1.100",
			RemotePort: 3306,
			Network:    "tcp",
		}})
		if err := v.WriteConfig(); err != nil {
			panic(err)
		}
		fmt.Print("示例已生成:./demo_proxy." + *t)
		return
	} else {
		//本地监听 端口
		go handler.Route(*p)
		var wg sync.WaitGroup
		c, _ := filepath.Abs(*c)
		v.SetConfigFile(c)
		ext := path.Ext(c)[1:]
		if ext == "" {
			ext = *t
		}
		v.SetConfigType(ext)
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		config := handler.Config{}
		if err := v.Unmarshal(&config); err != nil {
			panic(err)
		}
		for _, proxyConfig := range config.ProxyConfig {
			log.Println(proxyConfig)
			wg.Add(1)
			if proxyConfig.Enable {
				time.Sleep(time.Millisecond * 50)
				go handler.ProxyHandler(proxyConfig, &wg, &config.JhChannel)
			}
		}
		wg.Wait()
		fmt.Println("所有进程结束")
	}
}