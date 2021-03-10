package handler

type Config struct {
	ProxyConfig []ProxyConfig
	JhChannel   JHChannel
}

type ProxyConfig struct {
	LocalPort  int
	RemoteIp   string
	RemotePort int
	Enable     bool
	Network    string
}

type JHChannel struct {
	Enable bool
	Url    string
}
