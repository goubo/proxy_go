package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"path"
	"path/filepath"
	"sync"
)

// 读取配置文件
// 配置文件格式
// 解析出所有的配置
// 读取所有本地端口
// 读取所有本地端口对应远端端口
// 监听所有接口
// 创建代理到指定端口

var t = flag.String("t", "yaml", "配置文件格式, 支持 json|yaml, 自动读取文件后缀,无后缀需要手动指定")
var c = flag.String("c", "./demo_config.yaml", "指定配置文件")
var g = flag.Bool("g", false, "在当前目录生成示例配置文件")

type Config struct {
	ProxyConfig []ProxyConfig
}

type ProxyConfig struct {
	LocalPort  int
	RemoteIp   string
	RemotePort int
	Enable     bool
	Network    string
}

func ProxyHandler(conf ProxyConfig, wg *sync.WaitGroup) {
	listener, err := net.Listen(conf.Network, fmt.Sprintf(":%d", conf.LocalPort))
	if err != nil {
		log.Fatalf("端口监听失败 %v\n", err)
	}
	defer wg.Done()
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("访问链接获取失败", err)
		}
		dConn, err := net.Dial(conf.Network, fmt.Sprintf("%s:%d", conf.RemoteIp, conf.RemotePort))
		if err != nil {
			log.Println("创建连接失败", err)
		}
		go io.Copy(conn, dConn)
		go io.Copy(dConn, conn)
	}
}

func main() {
	flag.Parse()
	v := viper.New()
	if *g {
		v.SetConfigFile("./demo_proxy." + *t)
		v.Set("ProxyConfig", []ProxyConfig{{
			LocalPort:  43306,
			RemoteIp:   "192.168.1.100",
			RemotePort: 3306,
			Network:    "tcp",
		}})
		if err := v.WriteConfig(); err != nil {
			panic(err)
		}
		fmt.Print("示例已生成:./demo_config." + *t)
	} else {
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
		config := Config{}
		if err := v.Unmarshal(&config); err != nil {
			panic(err)
		}
		for _, proxyConfig := range config.ProxyConfig {
			fmt.Println(proxyConfig)
			wg.Add(1)
			if proxyConfig.Enable {
				go ProxyHandler(proxyConfig, &wg)
			}
		}
		wg.Wait()
		fmt.Println("所有进程结束")
	}
}
