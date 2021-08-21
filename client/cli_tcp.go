package client

import (
	"net"
	"strings"
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

	util.DisplaySocketOption(conn)

	var sum int
	var count int

	msg := "h"
	msg = strings.Repeat(msg, 69411)
	for {
		util.Notice("start write")
		timeout := time.Now().Add(15 * time.Second)
		conn.SetWriteDeadline(timeout)
		packetSize, err := conn.Write([]byte(msg))
		if err != nil {
			util.Notice("write failed, count:%d, error_msg:%s", count, err.Error())
			return err
		}
		sum += packetSize
		count += 1
		util.Notice("write packetSize:%d, count:%d, sum:%d", packetSize, count, sum)
		//return nil
	}
	return nil
}
