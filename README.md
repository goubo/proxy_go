# proxy_go

golang端口代理


-g 生成配置文件模板

-c 指定配置文件目录

代理远端ip:port到本机指定port

```yaml
proxyconfig:
- localport: 43306
  remoteip: 192.168.1.100
  remoteport: 3306
  enable: false
  network: tcp
```
