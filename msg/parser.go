package msg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
)

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

	len := binary.BigEndian.Uint32(buf)
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
	binary.Write(&buf, binary.BigEndian, uint32(len))

	_, err = buf.Write(msgByte)
	return err
}
