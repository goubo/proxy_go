package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	handler "proxy/hander"
	"sync"
)

// 读取配置文件
// 配置文件格式
// 解析出所有的配置
// 读取所有本地端口
// 读取所有本地端口对应远端端口
// 监听所有接口
// 创建代理到指定端口

var c = flag.String("c", "./demo_proxy.yaml", "指定配置文件")
var g = flag.Bool("g", false, "在当前目录生成示例配置文件")

func main() {
	flag.Parse()
	v := viper.New()
	c, _ := filepath.Abs(*c)
	var wg sync.WaitGroup

	v.SetConfigFile(c)
	v.SetConfigType("yaml")

	if *g {
		v.Set("ProxyConfig", []handler.ProxyConfig{{
			LocalPort:  43306,
			RemoteIp:   "192.168.1.100",
			RemotePort: 3306,
			Network:    "tcp",
		}})
		if err := v.WriteConfig(); err != nil {
			panic(err)
		}
		fmt.Print("示例已生成:", c)

	} else {
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
		config := handler.Config{}
		if err := v.Unmarshal(&config); err != nil {
			panic(err)
		}
		for _, proxyConfig := range config.ProxyConfig {
			fmt.Println(proxyConfig)
			wg.Add(1)
			if proxyConfig.Enable {
				go handler.ProxyHandler(proxyConfig, &wg)
			}
		}
	}
	wg.Wait()
	fmt.Println("所有进程结束")
}
