package server

import (
	"context"
	"log"

	"github.com/glstr/gwatcher/msg"
	"github.com/glstr/gwatcher/util"
	quic "github.com/lucas-clemente/quic-go"
)

type QuicServer struct {
	addr    string
	cancelF context.CancelFunc
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
		ctx, cancelF := context.WithCancel(context.Background())
		s.cancelF = cancelF
		sess, err := listener.Accept(ctx)
		if err != nil {
			util.Notice("quic accept failed, error_msg:%s", err.Error())
			continue
		}

		//go handleSession(sess)
		util.Notice("new sess accept")
		handler := NewQuicHanlder(sess)
		handler.Start()
	}
}

func (s *QuicServer) Stop() {
	if s.cancelF != nil {
		s.cancelF()
	}
}

//func handleSession(sess quic.Session) error {
//	stream, err := sess.AcceptStream(context.Background())
//	if err != nil {
//		log.Printf("accept stream failed, error_msg:%s", err.Error())
//		return err
//	}
//
//	log.Printf("stream_id:%d", stream.StreamID())
//	// Echo through the loggingWriter
//	return err
//}

type QuicHandler struct {
	s          quic.Session
	readQueue  chan *msg.MessageContainer
	writeQueue chan *msg.MessageContainer
	done       <-chan struct{}
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

func (qh *QuicHandler) Stop() error {
	err := qh.s.CloseWithError(quic.ApplicationErrorCode(0), "")
	if err != nil {
		util.Notice("close session failed:%s", err.Error())
	}
	return err
}

func (qh *QuicHandler) ReadLoop() error {
	for {
		select {
		case <-qh.done:
			util.Notice("read loop exit")
		default:
		}
		stream, err := qh.s.AcceptStream(context.Background())
		if err != nil {
			util.Notice("accept stream failed:%s", err.Error())
			return err
		}

		util.Notice("new stream:%d", stream.StreamID())
		go func() {
			err := qh.readFromStream(stream)
			if err != nil {
				util.Notice("read from stream failed:%s", err.Error())
				qh.Stop()
				return
			}
		}()
	}
}

func (qh *QuicHandler) readFromStream(stream quic.Stream) error {
	defer stream.CancelRead(0)
	parser := msg.NewParser()
	msg, err := parser.Unmarshal(stream)
	if err != nil {
		util.Notice("parse stream packet failed:%s", err.Error())
		return err
	}

	util.Notice("server read:%v", msg)
	qh.readQueue <- msg
	return nil
}

func (qh *QuicHandler) Process() {
	for {
		select {
		case <-qh.done:
			util.Notice("process done")
			return
		case req := <-qh.readQueue:
			res := qh.process(req)
			qh.writeQueue <- res
		}
	}
}

func (qh *QuicHandler) WriteLoop() {
	for {
		select {
		case <-qh.done:
			util.Notice("write loop exit")
			return
		case res := <-qh.writeQueue:
			err := qh.WriteToStream(res)
			if err != nil {
				util.Notice("write to stream failed:%s", err.Error())
				qh.Stop()
			}
		}
	}
}

func (qh *QuicHandler) WriteToStream(msgContainer *msg.MessageContainer) error {
	stream, err := qh.s.OpenStream()
	if err != nil {
		util.Notice("open stream failed")
		return err
	}
	defer stream.Close()
	parser := msg.NewParser()
	err = parser.Marshal(stream, msgContainer)
	if err != nil {
		util.Notice("parser marshal failed")
		return err
	}
	util.Notice("server write:%v", msgContainer)
	return nil
}

func (qh *QuicHandler) process(req *msg.MessageContainer) *msg.MessageContainer {
	res := &msg.MessageContainer{}
	res.Id = req.Id + 100000
	res.Data = req.Data
	return res
}

type Processor struct{}

func (p *Processor) Process(req *msg.MessageContainer) *msg.MessageContainer {
	return nil
}
