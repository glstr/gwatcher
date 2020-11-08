package server

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"time"

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

		//go s.handleRequest(conn, doNothing)
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

func echo(conn net.Conn, done <-chan struct{}) error {
	for {
		select {
		case <-done:
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

// rec eof and send req
func reqAfterEOF(conn net.Conn, done <-chan struct{}) error {
	for {
		buf := make([]byte, 1<<10)
		_, err := conn.Read(buf)
		if err != nil {
			util.Notice("read failed, error_msg:%s", err.Error())
			if err == io.EOF {
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				time.Sleep(2 * time.Second)
				err := conn.Close()
				if err != nil {
					util.Notice("close failed, error_msg:%s", err.Error())
					return err
				}

				for {
					count, err := conn.Write([]byte("hello world"))
					if err != nil {
						util.Notice("eof write failed, error_msg:%s", err.Error())
						return err
					}
					util.Notice("write coutn:%d", count)
				}

			}
			return err
		}
	}
}

// rec only once and do nothing
func doNothing(conn net.Conn, done <-chan struct{}) error {
	needRead := true
	for {
		select {
		case <-done:
			return nil
		default:
			if needRead {
				buf := make([]byte, 100)
				_, err := conn.Read(buf)
				if err != nil {
					util.Notice("read failed, error_msg:%s", err.Error())
					return err
				}
				needRead = false
			} else {
				time.Sleep(10 * time.Second)
			}

		}
	}
	return nil
}

func (s *TlsServer) Stop() {
	close(s.done)
}
