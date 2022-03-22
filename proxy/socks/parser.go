package socks

import (
	"io"
	"sync"

	"github.com/glstr/gwatcher/util"
)

type Parser interface {
	Write(input []byte) (int, error)
	Read(input []byte) (int, error)
	Close() error
}

func NewParser(localIP, remoteIP string, rwc io.ReadWriteCloser) Parser {
	parser := &DefaultParser{
		localIP:     localIP,
		remoteIP:    remoteIP,
		rwc:         rwc,
		readStream:  make(chan []byte, 1024*1024),
		writeStream: make(chan []byte, 1024*1024),
		parserFunc:  DefaultParseFunc,
	}

	go parser.parseReadStream()
	go parser.parseWriteStream()
	return parser
}

func DefaultParseFunc(dataStream <-chan []byte) error {
	for data := range dataStream {
		util.Notice("data:%d", len(data))
	}
	return nil
}

type DefaultParser struct {
	rwc io.ReadWriteCloser

	remoteIP string
	localIP  string

	readLen  int
	writeLen int

	readStream  chan []byte
	writeStream chan []byte
	closeOnce   sync.Once

	parserFunc func(<-chan []byte) error
}

func (p *DefaultParser) parseReadStream() error {
	return p.parserFunc(p.readStream)
}

func (p *DefaultParser) parseWriteStream() error {
	return p.parserFunc(p.writeStream)
}

func (p *DefaultParser) Write(input []byte) (int, error) {
	p.writeLen += len(input)
	p.writeStream <- input
	util.Notice("remote_ip:%s, write len:%d, sum:%d, body:%d",
		p.remoteIP,
		len(input),
		p.writeLen,
		len(input))
	return p.rwc.Write(input)
}

func (p *DefaultParser) Close() error {
	p.release()
	return p.rwc.Close()
}

func (p *DefaultParser) release() {
	p.closeOnce.Do(func() {
		close(p.readStream)
		close(p.writeStream)
	})
}

func (p *DefaultParser) Read(input []byte) (n int, err error) {
	p.readLen += len(input)
	p.readStream <- input
	util.Notice("remote_ip:%s, read len:%d, sum:%d, body:%d",
		p.remoteIP,
		len(input),
		p.readLen,
		len(input))
	return p.rwc.Read(input)
}
