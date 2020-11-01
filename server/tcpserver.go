package server

import (
	"bufio"
	"crypto/tls"
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

		go s.handleRequest(conn)
	}
	return nil
}

func (s *TlsServer) handleRequest(conn net.Conn) error {
	for {
		select {
		case <-s.done:
			return nil
		default:
		}
		buf := make([]byte, 100)
		_, err := conn.Read(buf)
		if err != nil {
			util.Notice("read failed, error_msg:%s", err.Error())
			return err
		}

		util.Notice("get content:%s", string(buf))
		writer := bufio.NewWriter(conn)
		count, err := writer.Write(buf)
		if err != nil {
			util.Notice("write failed, count:%d, error_msg:%s", count, err.Error())
			return err
		}
		writer.Flush()
	}
	return nil
}

func (s *TlsServer) Stop() {
	close(s.done)
}
