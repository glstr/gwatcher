package proxy

import (
	"crypto/tls"
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
	auther := NewAuther()
	err := auther.Auth(conn)
	if err != nil {
		util.Notice("auth failed, error:%v", err)
		return err
	}

	forwarder := NewForwarder(conn)
	return forwarder.Forward()
}

func SendSocks5Reply(conn net.Conn) error {
	reply := NewSocksReply()
	_, err := conn.Write(reply.Bytes())
	if err != nil {
		return errors.New("write rsp: " + err.Error())
	}
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

type Auther struct{}

func NewAuther() *Auther {
	return &Auther{}
}

func (a *Auther) Auth(conn net.Conn) error {
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

type Forwarder struct {
	srcConn net.Conn
	dstConn net.Conn
}

func NewForwarder(conn net.Conn) *Forwarder {
	return &Forwarder{
		srcConn: conn,
	}
}

func (f *Forwarder) Forward() error {
	err := f.connectDst()
	if err != nil {
		return err
	}

	return f.forward()
}

func (f *Forwarder) connectDst() error {
	var req Socks5Request
	err := req.Decode(f.srcConn)
	if err != nil {
		return err
	}

	switch req.Cmd {
	case CmdConnect:
		addr := net.JoinHostPort(req.Addr, fmt.Sprintf("%d", req.Port))
		//dstConn, err := net.Dial("tcp", net.JoinHostPort(req.Addr, fmt.Sprintf("%d", req.Port)))
		dstConn, err := tls.Dial("tcp", addr, nil)
		if err != nil {
			return err
		}
		err = SendSocks5Reply(f.srcConn)
		if err != nil {
			dstConn.Close()
			return err
		}
		f.dstConn = dstConn
		return nil

	default:
		return ErrNotSupportCmd
	}
}

func (f *Forwarder) forward() error {
	//conn := tls.Server(f.srcConn, util.GenerateTLSConfig())
	tlsConfigMaker := util.NewTlsConfigMaker()
	tlsConfig, err := tlsConfigMaker.MakeTls2Config()
	if err != nil {
		util.Notice("make tls config err:%v", err)
		return err
	}
	conn := tls.Server(f.srcConn, tlsConfig)
	dstParser := NewParser("dst"+f.dstConn.LocalAddr().String(), f.dstConn)
	//srcParser := NewParser("src"+f.srcConn.RemoteAddr().String(), f.srcConn)
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
