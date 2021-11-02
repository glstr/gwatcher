package client

import (
	"log"
	"net"
	"time"
)

type UdpClient struct {
	addr string
}

func NewUdpClient(addr string) *UdpClient {
	return &UdpClient{
		addr: addr,
	}
}

func (c *UdpClient) Start() error {
	conn, err := net.Dial("udp", c.addr)
	if err != nil {
		log.Printf("dial failed, error_msg:%s", err.Error())
		return err
	}

	word := "hello world"
	for {
		count, err := conn.Write([]byte(word))
		if err != nil {
			log.Printf("write failed, error_msg:%s", err.Error())
			return err
		}
		log.Printf("count:%d", count)
		time.Sleep(1 * time.Second)
	}
}
