package server

import (
	"context"
	"io"
	"log"

	"github.com/glstr/gwatcher/util"
	quic "github.com/lucas-clemente/quic-go"
)

type QuicServer struct {
	addr string
}

func NewQuicServer(addr string) *QuicServer {
	return &QuicServer{
		addr: addr,
	}
}

func (s *QuicServer) Start() error {
	log.Printf("start quic server")
	listener, err := quic.ListenAddr(s.addr, util.GenerateTLSConfig(), nil)
	if err != nil {
		return err
	}

	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("accept failed, error_msg:%s", err.Error())
			return err
		}

		log.Printf("accept success")
		go handleSession(sess)
	}
}

func (s *QuicServer) Stop() {
	return
}

func handleSession(sess quic.Session) error {
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		log.Printf("accept stream failed, error_msg:%s", err.Error())
		return err
	}

	log.Printf("stream_id:%d", stream.StreamID())
	// Echo through the loggingWriter
	_, err = io.Copy(util.LoggingWriter{stream}, stream)
	return err
}
