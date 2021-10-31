package protocol

type ProtocolType string

const (
	PQUIC ProtocolType = "quic"
	PUDP  ProtocolType = "udp"
)
