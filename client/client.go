package client

const (
	PQuic      = "quic"
	PQuicEarly = "quic_e"
	PUdp       = "udp"
	PTcp       = "tcp"
	PTls       = "tls"
	PHttp3     = "http3"
	PHttp      = "http"
)

type Client interface {
	Start() error
}

func NewClient(protocol, addr string) Client {
	switch protocol {
	case PUdp:
		return NewUdpClient(addr)
	case PQuic:
		return NewQuicClient(addr)
	case PQuicEarly:
		return NewQuicEarlyClient(addr)
	case PTcp:
		return NewTcpClient(addr)
	case PTls:
		return NewTlsClient(addr)
	case PHttp3:
		return NewHttp3Client(addr)
	case PHttp:
		return NewHttpClient(addr)
	default:
		return NewQuicClient(addr)
	}
}
