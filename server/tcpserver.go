package server

import (
	"errors"
	"net"

	"github.com/glstr/gwatcher/util"
)

type TcpServer struct {
	addr string
	done chan struct{}
}

func NewTcpServer(addr string) *TcpServer {
	return &TcpServer{
		addr: addr,
		done: make(chan struct{}),
	}
}

func (s *TcpServer) Start() error {
	util.Notice("start tcp server, addr:%s", s.addr)
	listener, err := net.Listen("tcp", s.addr)
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

		//go s.handleRequest(conn, reqAfterEOF)
		util.DisplaySocketOption(conn)
		go s.handleRequest(conn, doNothing)
	}

	return nil
}

func (s *TcpServer) handleRequest(conn net.Conn, handler ServerHandler) error {
	if handler == nil {
		return errors.New("invalid handler")
	}

	return handler(conn, s.done)
}

func (s *TcpServer) Stop() {
	close(s.done)
}
