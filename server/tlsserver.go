package server

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/glstr/gwatcher/util"
)

type TlsServer struct {
	addr string
	done chan struct{}
}

func NewTlsServer(addr string) *TlsServer {
	return &TlsServer{
		addr: addr,
		done: make(chan struct{}),
	}
}

func (s *TlsServer) Start() error {
	util.Notice("start tls server, addr:%s", s.addr)
	cerPath := "./conf/cert/glstr.cer"
	keyPath := "./conf/cert/server.key"
	cert, err := tls.LoadX509KeyPair(cerPath, keyPath)
	if err != nil {
		util.Notice("load tls key failed, error_msg:%s", err.Error())
		return err
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", s.addr, cfg)
	if err != nil {
		util.Notice("listen failed, error_msg:%s", err.Error())
		return err
	}

	for {
		conn, err := listener.Accept()
		util.Notice("accept conn")
		if err != nil {
			util.Notice("accept failed, error_msg:%s", err.Error())
			continue
		}

		go s.handleRequest(conn, reqAfterEOF)
	}
	return nil
}

type ServerHandler func(net.Conn, <-chan struct{}) error

func (s *TlsServer) handleRequest(conn net.Conn, handler ServerHandler) error {
	if handler == nil {
		return errors.New("invalid handler")
	}

	return handler(conn, s.done)
}

func (s *TlsServer) Stop() {
	close(s.done)
}
