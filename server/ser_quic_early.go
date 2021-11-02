package server

import (
	"context"
	"io"
	"log"

	"github.com/glstr/gwatcher/util"
	"github.com/lucas-clemente/quic-go"
)

type QuicEarlyServer struct {
	addr string
}

func NewQuicEarlyServer(addr string) *QuicEarlyServer {
	return &QuicEarlyServer{
		addr: addr,
	}
}

func (s *QuicEarlyServer) Start() error {
	util.Notice("start quic early server")

	earlyListener, err := quic.ListenAddrEarly(s.addr, util.GenerateTLSConfig(), nil)
	if err != nil {
		return err
	}

	for {
		earlySess, err := earlyListener.Accept(context.Background())
		if err != nil {
			continue
		}

		util.Notice("get early session")
		go handleEarlySession(earlySess)
	}
}

func (s *QuicEarlyServer) Stop() {
}

func handleEarlySession(sess quic.EarlySession) error {
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
