package client

import (
	"crypto/tls"
	"time"

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
		MaxVersion:         tls.VersionTLS12,
	}
	conn, err := tls.Dial("tcp", c.addr, cfg)
	if err != nil {
		util.Notice("dial tls failed, error_msg:%s", err.Error())
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
		//writer := bufio.NewWriter(conn)
		util.Notice("start write")
		timeout := time.Now().Add(15 * time.Second)
		conn.SetWriteDeadline(timeout)
		count, err := conn.Write([]byte(msg))
		if err != nil {
			util.Notice("write failed, count:%d, error_msg:%s", count, err.Error())
			return err
		}
		time.Sleep(3 * time.Second)
	}
}
