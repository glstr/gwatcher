package proxy

import (
	"errors"

	"github.com/glstr/gwatcher/proxy/socks"
)

var (
	ErrNotSupportProxyType = errors.New("not support proxy type")
)

type Proxy interface {
	Start() error
	Stop() error
}

type ProxyConfig struct {
	Host         string
	Port         string
	FileHostPort string
}

type ProxyType int

const (
	TypeSocks5Proxy ProxyType = iota
)

type makeProxyFunc func(config *ProxyConfig) Proxy

var ProxyType2ProxyMaker map[ProxyType]makeProxyFunc = map[ProxyType]makeProxyFunc{
	TypeSocks5Proxy: NewSocks5Proxy,
}

func NewProxy(t ProxyType, config *ProxyConfig) (Proxy, error) {
	if f, ok := ProxyType2ProxyMaker[t]; ok {
		return f(config), nil
	}

	return nil, ErrNotSupportProxyType
}

func NewSocks5Proxy(config *ProxyConfig) Proxy {
	c := &socks.Socks5ProxyConfig{
		Host:         config.Host,
		Port:         config.Port,
		FileHostPort: config.FileHostPort,
	}
	return socks.NewSocks5Proxy(c)
}
