package proxy

import "errors"

var (
	ErrNotSupportProxyType = errors.New("not support proxy type")
	ErrVersionFailed       = errors.New("socks5 version err")
	ErrNotSupportCmd       = errors.New("not support cmd")
	ErrMessageInvalid      = errors.New("socks5 message invalid")
)

type Proxy interface {
	Start() error
	Stop() error
}

type ProxyConfig struct {
	Host string
	Port string
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
