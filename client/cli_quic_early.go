package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"github.com/glstr/gwatcher/util"
	"github.com/lucas-clemente/quic-go"
)

type QuicEarlyClient struct {
	addr string
}

func NewQuicEarlyClient(addr string) *QuicEarlyClient {
	return &QuicEarlyClient{addr: addr}
}

func (c *QuicEarlyClient) Start() error {
	util.Notice("start quic early client")
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	session, err := quic.DialAddrEarly(c.addr, tlsConf, nil)
	if err != nil {
		return err
	}

	if err != nil {
		log.Printf("dial failed, error_msg:%s", err.Error())
		return err
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("open stream failed, error_msg:%s", err.Error())
		return err
	}

	defer stream.Close()
	message := "hello world1"
	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Printf("write stream failed, error_msg:%s", err.Error())
		return err
	}

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		log.Printf("read res failed, error_msg:%s", err.Error())
		return err
	}
	fmt.Printf("Client: Got '%s'\n", buf)

	return nil
}
