package socks

import (
	"errors"
	"net"

	"github.com/glstr/gwatcher/util"
)

type Socks5Proxy struct {
	Host string
	Port string
	done chan struct{}
}

func NewSocks5Proxy(host string, port string) *Socks5Proxy {
	return &Socks5Proxy{
		Host: host,
		Port: port,
		done: make(chan struct{}),
	}
}

func (p *Socks5Proxy) Start() error {
	address := net.JoinHostPort(p.Host, p.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go p.handleTcpListener(listener)
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
		util.Notice("get conn")
		go p.handleConn(conn)
	}
}

func (p *Socks5Proxy) handleConn(conn net.Conn) error {
	initHandler := NewInitMsgHandler()
	err := initHandler.Handle(conn)
	if err != nil {
		util.Notice("auth failed, error:%v", err)
		return err
	}

	handshakeHandler := NewHandshakeHandler()
	info, err := handshakeHandler.Handle(conn)
	if err != nil {
		util.Notice("handshake failed, error:%v", err)
		return err
	}

	relaySer := NewTcpRelayServer(conn, info.Addr)
	return relaySer.Relay()
}

func SendSocks5Reply(conn net.Conn) error {
	reply := NewSocksReply()
	_, err := conn.Write(reply.Bytes())
	if err != nil {
		return errors.New("write rsp: " + err.Error())
	}
	return nil
}

func (p *Socks5Proxy) waitStop() {
	<-p.done
}

func (p *Socks5Proxy) Stop() error {
	close(p.done)
	return nil
}
