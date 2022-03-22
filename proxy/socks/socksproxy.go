package socks

import (
	"errors"
	"net"
	"net/http"

	"github.com/glstr/gwatcher/util"
)

type Socks5ProxyConfig struct {
	Host         string
	Port         string
	FileHostPort string
}

type Socks5Proxy struct {
	Host         string
	Port         string
	FileHostPort string
	done         chan struct{}
}

func NewSocks5Proxy(c *Socks5ProxyConfig) *Socks5Proxy {
	return &Socks5Proxy{
		Host:         c.Host,
		Port:         c.Port,
		FileHostPort: c.FileHostPort,
		done:         make(chan struct{}),
	}
}

func (p *Socks5Proxy) startFileServer() error {
	util.Notice("start file server\n sockurl: http://127.0.0.1:8886/sock.pac\n caurl: http://127.0.0.1:8886/ca.crt")
	err := http.ListenAndServe(":"+p.FileHostPort, http.FileServer(http.Dir("./static")))
	if err != nil {
		util.Notice("start file server failed:%s", err.Error())
		return err
	}

	return nil
}

func (p *Socks5Proxy) Start() error {
	go p.startFileServer()

	address := net.JoinHostPort(p.Host, p.Port)
	util.Notice("start socket server:%s", address)
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

	relaySer := NewRelayServer(info)
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
