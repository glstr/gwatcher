package server

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/glstr/gwatcher/util"
)

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
	util.DisplaySocketOption(conn)
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
