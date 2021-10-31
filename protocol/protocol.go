package protocol

import "io"

type Message interface {
	Display() string
}

type Parser interface {
	Parse(io.Reader) (Message, error)
}

var (
	ParserMap = map[ProtocolType]Parser{
		PQUIC: &QuicParser{},
	}
)

func GetParser(p ProtocolType) (Parser, error) {
	if parser, ok := ParserMap[p]; ok {
		return parser, nil
	}

	return nil, ErrProtocolNotSupport
}
