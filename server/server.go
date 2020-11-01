package server

const (
	PQuic  = "quic"
	PUdp   = "udp"
	PHttp3 = "http3"
	PTls   = "tls"
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
	case PTls:
		return NewTlsServer(addr)
	default:
		return NewQuicServer(addr)
	}
}
