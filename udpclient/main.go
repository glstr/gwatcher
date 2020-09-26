package main

import (
	"log"
	"net"
	"time"
)

func main() {
	addr := "127.0.0.1:8888"
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Printf("dial failed, error_msg:%s", err.Error())
		return
	}

	word := "hello world"
	for {
		count, err := conn.Write([]byte(word))
		if err != nil {
			log.Printf("write failed, error_msg:%s", err.Error())
			return
		}
		log.Printf("count:%d", count)
		time.Sleep(1 * time.Second)
	}
}
