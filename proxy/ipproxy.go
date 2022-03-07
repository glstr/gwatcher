package proxy

import (
	"net"

	"github.com/glstr/gwatcher/util"
)

type IPProxyer struct {
	Addr string

	// forward for input and output ip
	AddrsMap map[string]string
	Done     chan struct{}
}

func (p *IPProxyer) Start() error {
	for {
		select {
		case <-p.Done:
			return nil
		default:
		}
		conn, err := net.ListenPacket("ip", "127.0.0.1")
		if err != nil {
			util.Notice("listen packet conn failed, error_msg:%s", err.Error())
			return err
		}

		go p.handlePacketConn(conn)
	}
}

func (p *IPProxyer) handlePacketConn(conn net.PacketConn) error {
	//packetBuffer := buffer.GetPacketBuffer()
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(buffer)
	if err != nil {
		util.Notice("read from failed, error_msg:%s", err.Error())
	}

	util.Notice("from addr:%v, data len:%d", addr, n)
	return nil
}

func (p *IPProxyer) Stop() error {
	close(p.Done)
	return nil
}
