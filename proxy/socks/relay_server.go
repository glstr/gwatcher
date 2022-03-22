package socks

import (
	"crypto/tls"
	"io"
	"net"

	"github.com/glstr/gwatcher/util"
)

type RelayServer interface {
	Relay() error
}

type RelayInfo struct {
	SrcConn net.Conn
	Cmd     byte
	Addr    string
}

func NewRelayServer(info RelayInfo) RelayServer {
	switch info.Cmd {
	case CmdBind, CmdConnect:
		return &TcpRelayServer{
			srcConn: info.SrcConn,
			dstAddr: info.Addr,
		}
	case CmdUdpASSOCIATE:
		return &UdpRelayServer{}
	default:
		return &TcpRelayServer{
			srcConn: info.SrcConn,
			dstAddr: info.Addr,
		}
	}

}

type TcpRelayServer struct {
	srcConn net.Conn
	dstAddr string
}

func (s *TcpRelayServer) Relay() error {
	dstConn, err := tls.Dial("tcp", s.dstAddr, nil)
	if err != nil {
		return err
	}

	tlsConfigMaker := util.NewTlsConfigMaker()
	tlsConfig, err := tlsConfigMaker.MakeTls2Config()
	if err != nil {
		util.Notice("make tls config err:%v", err)
		return err
	}
	conn := tls.Server(s.srcConn, tlsConfig)
	dstParser := NewParser(dstConn.LocalAddr().String(), dstConn.RemoteAddr().String(), dstConn)

	forward := func(src, dest io.ReadWriteCloser) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}

	go forward(dstParser, conn)
	go forward(conn, dstParser)
	return nil

}

type UdpRelayServer struct{}

func (s *UdpRelayServer) Relay() error {
	return nil
}
