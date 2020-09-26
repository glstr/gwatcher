package main

import (
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8888")
	if err != nil {
		log.Printf("resolve addr failed, error_msg:%s", err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("listen conn failed, error_msg:%s", err.Error())
		return
	}

	done := make(chan struct{})
	for {
		select {
		case <-done:
			return
		default:
		}

		content := make([]byte, 1000)
		count, remoteAddr, err := conn.ReadFromUDP(content)
		if err != nil {
			log.Printf("read data failed, error_msg:%s", err.Error())
			return
		}

		//log.Printf("count:%d, addr:%v, data:%s", count, *remoteAddr, string(content))
		log.Printf("count:%d, addr:%v, data:%s", count, remoteAddr, string(content))
	}
}
