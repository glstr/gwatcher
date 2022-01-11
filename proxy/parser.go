package proxy

import (
	"io"

	"github.com/glstr/gwatcher/util"
)

type Parser struct {
	Name string

	rwc io.ReadWriteCloser

	readLen  int
	writeLen int
}

func NewParser(name string, rwc io.ReadWriteCloser) *Parser {
	return &Parser{
		Name: name,
		rwc:  rwc,
	}
}

func (p *Parser) Write(input []byte) (int, error) {
	p.writeLen += len(input)
	util.Notice("from:%s, write len:%d, sum:%d, body:%s",
		p.Name,
		len(input),
		p.writeLen,
		string(input))
	return p.rwc.Write(input)
}

func (p *Parser) Close() error {
	return p.rwc.Close()
}

func (p *Parser) Read(input []byte) (n int, err error) {
	p.readLen += len(input)
	util.Notice("from:%s, read len:%d, sum:%d, body:%s",
		p.Name,
		len(input),
		p.readLen,
		string(input))
	return p.rwc.Read(input)
}
