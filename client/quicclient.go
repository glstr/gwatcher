package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	quic "github.com/lucas-clemente/quic-go"
)

const (
	PQuic  = "quic"
	PUdp   = "udp"
	PTcp   = "tcp"
	PTls   = "tls"
	PHttp3 = "http3"
	PHttp  = "http"
)

type Client interface {
	Start() error
}

func NewClient(protocol, addr string) Client {
	switch protocol {
	case PUdp:
		return NewUdpClient(addr)
	case PQuic:
		return NewQuicClient(addr)
	case PTcp:
		return NewTcpClient(addr)
	case PTls:
		return NewTlsClient(addr)
	case PHttp3:
		return NewHttp3Client(addr)
	case PHttp:
		return NewHttpClient(addr)
	default:
		return NewQuicClient(addr)
	}
}

type QuicClient struct {
	addr string
}

func NewQuicClient(addr string) *QuicClient {
	return &QuicClient{
		addr: addr,
	}
}

func (c *QuicClient) Start() error {
	log.Printf("start quic client")
	tlsConf := &tls.Config{
		InsecureSkipVerify: false,
		NextProtos:         []string{"lcp"},
	}
	session, err := quic.DialAddr(c.addr, tlsConf, nil)
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
