package server

import "strings"

type Server interface {
	Start() error
	Stop()
}
type serverFunc func(protocol ProtocolType, handler HandlerType, addr string) Server

var serverMap = map[ProtocolType]serverFunc{
	PUdp:   GetUdpServer,
	PQuic:  GetQuicServer,
	PTcp:   GetTcpServer,
	PTls:   GetTlsServer,
	PHttp:  GetHttpServer,
	PHttp3: GetHttp3Server,
}

func DisplayProtocols() string {
	var res []string
	for protocol := range serverMap {
		res = append(res, string(protocol))
	}

	return strings.Join(res, ",")
}

func GetServer(protocol ProtocolType, htype HandlerType, addr string) (Server, error) {
	if serverFunc, ok := serverMap[protocol]; ok {
		return serverFunc(protocol, htype, addr), nil
	}
	return nil, ErrNoServer
}

func GetUdpServer(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewUdpServer(addr)
}

func GetQuicServer(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewQuicServer(addr)
}

func GetTcpServer(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewTcpServer(addr)
}

func GetTlsServer(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewTlsServer(htype, addr)
}

func GetHttpServer(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewHttpServer(addr)
}

func GetHttp3Server(protocol ProtocolType, htype HandlerType, addr string) Server {
	return NewHttp3Server(addr)
}
