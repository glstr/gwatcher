package client

import (
	"net"
	"time"

	"github.com/glstr/gwatcher/util"
)

type TcpClient struct {
	addr string
}

func NewTcpClient(addr string) *TcpClient {
	return &TcpClient{
		addr: addr,
	}
}

func (c *TcpClient) Start() error {
	util.Notice("start tcp client, add:%s", c.addr)
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		util.Notice("dial tcp failed, error_msg:%s", err.Error())
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			util.Notice("close failed, error_msg:%s", err.Error())
		}
	}()

	for {
		msg := "hello world"
		util.Notice("start write")
		timeout := time.Now().Add(15 * time.Second)
		conn.SetWriteDeadline(timeout)
		count, err := conn.Write([]byte(msg))
		if err != nil {
			util.Notice("write failed, count:%d, error_msg:%s", count, err.Error())
			return err
		}
		util.Notice("write count:%d", count)
		return nil
	}
	return nil
}
