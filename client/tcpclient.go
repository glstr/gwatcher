package client

import (
	"bufio"
	"crypto/tls"

	"github.com/glstr/gwatcher/util"
)

type TlsClient struct {
	addr string
}

func NewTlsClient(addr string) *TlsClient {
	return &TlsClient{
		addr: addr,
	}
}

func (c *TlsClient) Start() error {
	util.Notice("start tls client, add:%s", c.addr)
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", c.addr, cfg)
	if err != nil {
		util.Notice("dial failed, error_msg:%s", err.Error())
		return err
	}
	//defer conn.Close()
	msg := "hello world"
	writer := bufio.NewWriter(conn)
	count, err := writer.Write([]byte(msg))
	if err != nil {
		util.Notice("write failed, count:%d, error_msg:%s", count, err.Error())
		return err
	}
	util.Notice("write count:%d", count)
	return writer.Flush()
}
