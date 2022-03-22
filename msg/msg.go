package msg

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

type MsgType int

const (
	TypeRequest MsgType = iota
	TypeResponse
	TypePush
)

type MsgMetaExt struct {
	StatusCode int32
	Msg        string
}

type MsgMeta struct {
	Service   string
	Method    string
	Id        int64
	Timestamp int64
	Type      MsgType
}

type Msg struct {
	MagicNum byte
	DataLen  int
	MetaLen  int
	Meta     MsgMeta
	Payload  []byte

	// meta + payload
	data []byte
}

func NewMsg() *Msg {
	return &Msg{}
}

func (m *Msg) Decode(r io.Reader) error {
	headerBuffer := make([]byte, 5)
	_, err := io.ReadFull(r, headerBuffer)
	if err != nil {
		return err
	}

	m.MagicNum = byte(headerBuffer[0])
	if m.MagicNum != 'P' {
		return errors.New("protocl error")
	}

	m.DataLen = int(binary.BigEndian.Uint16(headerBuffer[1:3]))
	m.MetaLen = int(binary.BigEndian.Uint16(headerBuffer[3:5]))

	if len(m.data) < int(m.DataLen) {
		m.data = make([]byte, m.DataLen)
	}

	_, err = io.ReadFull(r, m.data)
	if err != nil {
		return err
	}

	metaBuffer := m.data[:m.MetaLen]
	err = json.Unmarshal(metaBuffer, &m.Meta)
	if err != nil {
		return err
	}

	m.Payload = m.data[m.MetaLen:]
	return nil
}

func (m *Msg) Encode() ([]byte, error) {
	metaBuffer, err := json.Marshal(m.Meta)
	if err != nil {
		return nil, err
	}

	m.MetaLen = len(metaBuffer)
	m.DataLen = len(m.Payload) + m.MetaLen

	buffer := make([]byte, 5+m.DataLen)
	buffer[0] = m.MagicNum
	binary.BigEndian.PutUint16(buffer[1:3], uint16(m.DataLen))
	binary.BigEndian.PutUint16(buffer[3:5], uint16(m.MetaLen))
	copy(buffer[5:5+m.MetaLen], metaBuffer)
	copy(buffer[5+m.MetaLen:], m.Payload)

	return buffer, nil
}
