# proxy_go

golang端口代理


-g 生成配置文件模板

-c 指定配置文件目录

代理远端ip:port到本机指定port

```yaml
proxyconfig:
- localport: 43306   # 本地端口
  remoteip: 192.168.1.100  # 远端ip
  remoteport: 3306   # 远端端口
  enable: false   # 是否启用
  network: tcp    # 协议类型
```
