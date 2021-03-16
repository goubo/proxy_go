package handler

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func ProxyHandler(conf ProxyConfig, wg *sync.WaitGroup, channel *JHChannel) {

	remoteIp, remotePort, err := getChannel(conf, channel)

	if err != nil {
		fmt.Printf("通道申请失败 %v !\n", err)
	}
	conf.RemotePort = remotePort
	conf.RemoteIp = remoteIp
	listener, err := net.Listen(conf.Network, fmt.Sprintf(":%d", conf.LocalPort))
	if err != nil {
		log.Fatalf("端口监听失败 %v\n", err)
		return
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

func getChannel(conf ProxyConfig, channel *JHChannel) (remoteIp string, remotePort int, err error) {

	if !channel.Enable {
		return conf.RemoteIp, conf.RemotePort, nil
	}

	return "", 0, nil

}
