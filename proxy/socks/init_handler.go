package socks

import (
	"errors"
	"fmt"
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

func (h *HandshakeHandler) Handle(conn net.Conn) (RelayInfo, error) {
	info := RelayInfo{
		SrcConn: conn,
	}

	var req Socks5Request
	err := req.Decode(conn)
	if err != nil {
		return info, err
	}

	info.Cmd = req.Cmd
	switch req.Cmd {
	case CmdConnect:
		addr := net.JoinHostPort(req.Addr, fmt.Sprintf("%d", req.Port))
		info.Addr = addr
		err = SendSocks5Reply(conn)
		return info, err
	case CmdUdpASSOCIATE:
		return info, err
	default:
		return info, ErrNotSupportCmd
	}
}
