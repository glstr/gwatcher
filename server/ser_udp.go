package server

import (
	"log"
	"net"
)

type Proccessor interface {
	Proccess(conn net.Conn) error
	Stop()
}

type EchoProccessor struct {
	done chan struct{}
}

func (p *EchoProccessor) Proccess(conn net.Conn) error {
	for {
		select {
		case <-p.done:
			log.Println("stop udp server")
			return nil
		default:
		}

		content := make([]byte, 1000)
		_, err := conn.Read(content)
		if err != nil {
			log.Printf("read data failed, error_msg:%s", err.Error())
			return err
		}
		log.Printf("content:%s", string(content))
	}
}

func (p *EchoProccessor) Stop() {
	p.done <- struct{}{}
}

type UdpServer struct {
	addr string
	p    Proccessor
}

func NewUdpServer(addr string) *UdpServer {
	return &UdpServer{
		addr: addr,
		p: &EchoProccessor{
			done: make(chan struct{}),
		},
	}
}

func (s *UdpServer) Start() error {
	log.Printf("start udp server")
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	return s.p.Proccess(conn)
}

func (s *UdpServer) Stop() {
	s.p.Stop()
}
