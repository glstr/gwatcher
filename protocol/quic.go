package protocol

import "io"

type QuicMessage struct{}

func (m *QuicMessage) Display() string {
	return "Hello world quic"
}

type QuicParser struct{}

func (p *QuicParser) Parse(r io.Reader) (Message, error) {
	return &QuicMessage{}, nil
}
