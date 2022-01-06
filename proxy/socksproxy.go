package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/glstr/gwatcher/util"
)

type Socks5Proxy struct {
	config ProxyConfig

	done chan struct{}
}

func NewSocks5Proxy(config *ProxyConfig) Proxy {
	return &Socks5Proxy{
		config: *config,
	}
}

func (p *Socks5Proxy) Start() error {
	address := net.JoinHostPort(p.config.Host, p.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go p.handleTcpListener(listener)

	packConn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	go p.handlePackConn(packConn)

	p.waitStop()
	return nil
}

func (p *Socks5Proxy) handleTcpListener(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			util.Notice("accpet failed, error_msg:%s", err.Error())
			return err
		}
		go p.handleConn(conn)
	}
}

func (p *Socks5Proxy) handleConn(conn net.Conn) error {
	if err := p.auth(conn); err != nil {
		util.Notice("auth failed, error_msg:%s", err.Error())
		return err
	}

	dstConn, err := p.connect(conn)
	if err != nil {
		util.Notice("connect failed, error_msg:%s", err.Error())
		return err
	}

	util.Notice("start forward")
	if err := p.forward(dstConn, conn); err != nil {
		util.Notice("forward failed, error_msg:%s", err.Error())
		return err
	}

	return nil
}

func (p *Socks5Proxy) auth(conn net.Conn) error {
	var auth Socks5Auth
	err := auth.Decode(conn)
	if err != nil {
		return err
	}

	if auth.Ver != VersionSocks5 {
		util.Notice("auth.Version:%d", auth.Ver)
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

func (p *Socks5Proxy) connect(conn net.Conn) (net.Conn, error) {
	var req Socks5Request
	err := req.Decode(conn)
	if err != nil {
		return nil, err
	}

	switch req.Cmd {
	case CmdConnect:
		dstConn, err := net.Dial("tcp", net.JoinHostPort(req.Addr, fmt.Sprintf("%d", req.Port)))
		if err != nil {
			return nil, err
		}
		err = p.sendReply(conn)
		if err != nil {
			dstConn.Close()
			return nil, err
		}
		return dstConn, nil

	default:
		return nil, ErrNotSupportCmd
	}
}

func (p *Socks5Proxy) sendReply(conn net.Conn) error {
	//	+----+-----+-------+------+----------+----------+
	//  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	//  +----+-----+-------+------+----------+----------+
	//  | 1  |  1  | X'00' |  1   | Variable |    2     |
	//  +----+-----+-------+------+----------+----------+
	_, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return errors.New("write rsp: " + err.Error())
	}
	return nil
}

func (p *Socks5Proxy) forward(dstConn, conn net.Conn) error {
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(dstConn, conn)
	go forward(conn, dstConn)
	return nil
}

func (p *Socks5Proxy) handlePackConn(packConn net.PacketConn) error {
	return nil
}

func (p *Socks5Proxy) waitStop() {
	<-p.done
}

func (p *Socks5Proxy) Stop() error {
	close(p.done)
	return nil
}
