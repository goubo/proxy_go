package handler

type Config struct {
	ProxyConfig []ProxyConfig
}

type ProxyConfig struct {
	LocalPort int
	RemoteIp string
	RemotePort int
	Enable bool
	Network string
}
