package socks

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/glstr/gwatcher/util"
)

type InitMsgHandler struct{}

func NewInitMsgHandler() *InitMsgHandler {
	return &InitMsgHandler{}
}

func (p *InitMsgHandler) Handle(conn net.Conn) error {
	var initMsg Socks5Init
	err := initMsg.Decode(conn)
	if err != nil {
		return err
	}

	if initMsg.Ver != VersionSocks5 {
		util.Notice("auth.Version:%d", initMsg.Ver)
		return errors.New("version err")
	}

	// return no need auth
	authReply := Socks5AuthReply{
		Ver:    VersionSocks5,
		Method: MethodNoAuth,
	}

	n, err := conn.Write(authReply.Encode())
	if err != nil || n != 2 {
		return errors.New("reply failed")
	}

	util.Notice("reply success")
	return nil
}

type HandshakeHandler struct{}

func NewHandshakeHandler() *HandshakeHandler {
	return &HandshakeHandler{}
}

type RelayInfo struct {
	Cmd  byte
	Addr string
}

func (h *HandshakeHandler) Handle(conn net.Conn) (RelayInfo, error) {
	var info RelayInfo

	var req Socks5Request
	err := req.Decode(conn)
	if err != nil {
		return info, err
	}

	switch req.Cmd {
	case CmdConnect:
		addr := net.JoinHostPort(req.Addr, fmt.Sprintf("%d", req.Port))
		info.Addr = addr
		err = SendSocks5Reply(conn)
		if err != nil {
			return info, err
		}
		return info, nil

	default:
		return info, ErrNotSupportCmd
	}
}

type TcpRelayServer struct {
	srcConn net.Conn

	dstAddr string
}

func NewTcpRelayServer(srcConn net.Conn, dstAddr string) *TcpRelayServer {
	return &TcpRelayServer{
		srcConn: srcConn,
		dstAddr: dstAddr,
	}
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
	dstParser := NewParser("dst"+dstConn.LocalAddr().String(), dstConn)
	srcParser := NewParser("src"+conn.RemoteAddr().String(), conn)

	forward := func(src, dest io.ReadWriteCloser) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}

	go forward(dstParser, srcParser)
	go forward(srcParser, dstParser)
	return nil

}
