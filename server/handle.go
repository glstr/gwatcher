package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"strings"
	"time"

	"github.com/glstr/gwatcher/util"
)

type ServerHandler func(net.Conn, <-chan struct{}) error

type Handler interface {
	handle(context context.Context, conn net.Conn) error
}

var handlers = map[HandlerType]Handler{
	HEcho: &EchoHandler{},
}

func DisplayHandlers() string {
	var res []string
	for htype := range handlers {
		res = append(res, string(htype))
	}

	return strings.Join(res, ",")
}

func GetHandler(htype HandlerType) (Handler, error) {
	if h, ok := handlers[htype]; ok {
		return h, nil
	}

	return nil, ErrNoHandler
}

type EchoHandler struct{}

func (h *EchoHandler) handle(ctx context.Context, conn net.Conn) error {
	for {
		select {
		case <-ctx.Done():
			util.Notice("handler exit")
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
}

// rec eof and send req
func reqAfterEOF(conn net.Conn, done <-chan struct{}) error {
	for {
		buf := make([]byte, 1<<10)
		_, err := conn.Read(buf)
		if err != nil {
			util.Notice("read failed, error_msg:%s", err.Error())
			if err == io.EOF {
				//conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				//time.Sleep(2 * time.Second)
				err := conn.Close()
				if err != nil {
					util.Notice("close failed, error_msg:%s", err.Error())
					return err
				}
				util.Notice("close success")

				for {
					count, err := conn.Write([]byte("hello world"))
					if err != nil {
						util.Notice("eof write failed, error_msg:%s", err.Error())

						//err = conn.Close()
						//if err != nil {
						//	util.Notice("close  failed, error_msg:%s", err.Error())
						//}
						//util.Notice("close success")
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
}
