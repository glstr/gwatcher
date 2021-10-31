package server

import (
	"context"
	"io"
	"log"

	"github.com/glstr/gwatcher/msg"
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

type QuicHandler struct {
	s          quic.Session
	readQueue  chan *msg.MessageContainer
	writeQueue chan *msg.MessageContainer
}

func NewQuicHanlder(s quic.Session) *QuicHandler {
	return &QuicHandler{
		s:          s,
		readQueue:  make(chan *msg.MessageContainer, 1024),
		writeQueue: make(chan *msg.MessageContainer, 1024),
	}
}

func (qh *QuicHandler) Start() error {
	go qh.ReadLoop()
	go qh.Process()
	go qh.WriteLoop()
	return nil
}

func (qh *QuicHandler) ReadLoop() error {
	for {
		stream, err := qh.s.AcceptStream(context.Background())
		if err != nil {
			log.Printf("accept stream failed:%s", err.Error())
			continue
		}
		go qh.readFromStream(stream)
	}
}

func (qh *QuicHandler) readFromStream(stream quic.Stream) error {
	defer stream.CancelRead(0)
	parser := msg.NewParser()
	msg, err := parser.Unmarshal(stream)
	if err != nil {
		log.Printf("parse stream packet failed:%s", err.Error())
		return err
	}

	qh.readQueue <- msg
	return nil
}

func (qh *QuicHandler) Process() {
	for {
		select {
		case req := <-qh.readQueue:
			res := qh.process(req)
			qh.writeQueue <- res
		}
	}
}

func (qh *QuicHandler) WriteLoop() {
	for {
		select {
		case res := <-qh.writeQueue:
			err := qh.WriteToStream(res)
			if err != nil {
				log.Printf("write failed:%s", err.Error())
			}
		}
	}
}

func (qh *QuicHandler) WriteToStream(res *msg.MessageContainer) error {
	stream, err := qh.s.OpenUniStream()
	if err != nil {
		log.Printf("open stream failed")
		return err
	}
	defer stream.CancelWrite(0)
	parser := msg.NewParser()
	err = parser.Marshal(stream, res)
	if err != nil {
		log.Printf("parser marshal failed")
		return err
	}
	return nil
}

func (qh *QuicHandler) process(req *msg.MessageContainer) *msg.MessageContainer {
	res := &msg.MessageContainer{}
	res.Id = req.Id + 100000
	res.Data = req.Data
	return res
}
