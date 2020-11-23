package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
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

type QuicHandler struct {
	s          quic.Session
	readQueue  chan *MessageContainer
	writeQueue chan *MessageContainer
}

func NewQuicHanlder(s quic.Session) *QuicHandler {
	return &QuicHandler{
		s:          s,
		readQueue:  make(chan *MessageContainer, 1024),
		writeQueue: make(chan *MessageContainer, 1024),
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
	parser := NewParser()
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

func (qh *QuicHandler) WriteToStream(res *MessageContainer) error {
	stream, err := qh.s.OpenUniStream()
	if err != nil {
		log.Printf("open stream failed")
		return err
	}
	defer stream.CancelWrite(0)
	parser := NewParser()
	err = parser.Marshal(stream, res)
	if err != nil {
		log.Printf("parser marshal failed")
		return err
	}
	return nil
}

func (qh *QuicHandler) process(req *MessageContainer) *MessageContainer {
	res := &MessageContainer{}
	res.Id = req.Id + 100000
	res.Data = req.Data
	return res
}

type MessageContainer struct {
	Id   uint64
	Data string
}

type Parser interface {
	Unmarshal(io.Reader) (*MessageContainer, error)
	Marshal(io.Writer, *MessageContainer) error
}

func NewParser() Parser {
	return &parser{}
}

// message packet format: len(4 bytes) + json body
type parser struct{}

func (p *parser) Unmarshal(reader io.Reader) (*MessageContainer, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		log.Printf("read len failed:%s", err.Error())
		return nil, err
	}

	len := binary.BigEndian.Uint64(buf)
	bodyBuf := make([]byte, len)
	_, err = io.ReadFull(reader, bodyBuf)
	if err != nil {
		log.Printf("read body failed:%s", err.Error())
		return nil, err
	}

	msg := &MessageContainer{}
	err = json.Unmarshal(bodyBuf, msg)
	if err != nil {
		log.Printf("get MessageContainer failed:%s", err.Error())
		return nil, err
	}

	return msg, nil
}

func (p *parser) Marshal(writer io.Writer, msg *MessageContainer) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Printf("json marshal failed:%s", err.Error())
		return err
	}

	len := len(msgByte)

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint64(len))

	_, err = buf.Write(msgByte)
	return err
}
