package server

const (
	PQuic  = "quic"
	PUdp   = "udp"
	PHttp3 = "http3"
	PTcp   = "tcp"
	PTls   = "tls"
	PHttp  = "http"
)

type Server interface {
	Start() error
	Stop()
}

func NewServer(protocol string, addr string) Server {
	switch protocol {
	case PUdp:
		return NewUdpServer(addr)
	case PQuic:
		return NewQuicServer(addr)
	case PHttp3:
		return NewHttp3Server(addr)
	case PTcp:
		return NewTcpServer(addr)
	case PTls:
		return NewTlsServer(addr)
	case PHttp:
		return NewHttpServer(addr)
	default:
		return NewQuicServer(addr)
	}
}
