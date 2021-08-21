package server

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/glstr/gwatcher/util"
)

type TlsServer struct {
	htype HandlerType
	addr  string

	ctx    context.Context
	cancel func()
}

func NewTlsServer(htype HandlerType, addr string) *TlsServer {
	ctx, f := context.WithCancel(context.Background())
	return &TlsServer{
		htype: htype,
		addr:  addr,

		ctx:    ctx,
		cancel: f,
	}
}

func (s *TlsServer) Start() error {
	util.Notice("start tls server, addr:%s", s.addr)
	cfg, err := GetTlsConfig()
	if err != nil {
		util.Notice("get tls config failed:%s", err.Error())
		return err
	}

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

		go s.handleRequest(conn)
	}
}

func (s *TlsServer) handleRequest(conn net.Conn) error {
	handler, err := GetHandler(s.htype)
	if err != nil {
		util.Notice("get handler failed, error_msg:%s", err.Error())
		return err
	}

	return handler.handle(s.ctx, conn)
}

func (s *TlsServer) Stop() {
	s.cancel()
}

func GetTlsConfig() (*tls.Config, error) {
	cerPath := "./conf/cert/glstr.cer"
	keyPath := "./conf/cert/server.key"
	cert, err := tls.LoadX509KeyPair(cerPath, keyPath)
	if err != nil {
		util.Notice("load tls key failed, error_msg:%s", err.Error())
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
