package server

type ProtocolType string

const (
	PQuic      ProtocolType = "quic"
	PQuicEarly ProtocolType = "quic_e"
	PUdp       ProtocolType = "udp"
	PHttp3     ProtocolType = "http3"
	PTcp       ProtocolType = "tcp"
	PTls       ProtocolType = "tls"
	PHttp      ProtocolType = "http"
)

type HandlerType string

const (
	HDothing HandlerType = "DoNothing"
	HEcho    HandlerType = "Echo"
)
