package protocol

import (
	"fmt"

	"github.com/glstr/gwatcher/util"
)

type CodecOption struct {
	FilePath string
	Protocol ProtocolType
}

type Codec struct{}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) Parse(option *CodecOption) error {
	parser, err := GetParser(option.Protocol)
	if err != nil {
		return err
	}

	reader, err := util.OpenReadFile(option.FilePath)
	if err != nil {
		return err
	}

	msg, err := parser.Parse(reader)
	if err != nil {
		return err
	}

	fmt.Printf("msg:%s done", msg.Display())
	return nil
}
