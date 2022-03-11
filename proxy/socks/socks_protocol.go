package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/glstr/gwatcher/buffer"
	"github.com/glstr/gwatcher/util"
)

const (
	VersionSocks5 = 0x05
)

const (
	CmdConnect      = 0x01
	CmdBind         = 0x02
	CmdUdpASSOCIATE = 0x03

	AtypIPV4   = 0x01
	AtypDomain = 0x03
	AtypIPV6   = 0x04
)

const (
	MethodNoAuth = 0x00
)

type Socks5Init struct {
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	Ver      int
	NMethods int
	Methods  []byte
}

func (s *Socks5Init) Decode(conn net.Conn) error {
	buffer := buffer.SmallBytesPool.Get().(*buffer.SmallBytes)
	defer buffer.Release()

	tmp := buffer.Data[:2]
	_, err := io.ReadFull(conn, tmp)
	if err != nil {
		return err
	}

	s.Ver, s.NMethods = int(tmp[0]), int(tmp[1])
	util.Notice("s.Ver:%d, s.NMethods:%d", s.Ver, s.NMethods)

	s.Methods = make([]byte, s.NMethods)
	_, err = io.ReadFull(conn, s.Methods)
	if err != nil {
		return err
	}

	return nil
}

type Socks5AuthReply struct {
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	Ver    byte
	Method byte
}

func (r *Socks5AuthReply) Encode() []byte {
	return []byte{r.Ver, r.Method}
}

type Socks5Request struct {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	Ver  byte
	Cmd  byte
	Rsv  byte
	Atyp byte
	Addr string
	Port uint16
}

func (req *Socks5Request) Decode(conn net.Conn) error {
	buffer := buffer.SmallBytesPool.Get().(*buffer.SmallBytes)
	defer buffer.Release()

	tmp := buffer.Data[:4]
	_, err := io.ReadFull(conn, tmp)
	if err != nil {
		return err
	}

	req.Ver, req.Cmd, req.Rsv, req.Atyp = tmp[0], tmp[1], tmp[2], tmp[3]
	if req.Ver != VersionSocks5 {
		return ErrVersionFailed
	}

	err = req.getAddr(conn)
	if err != nil {
		return err
	}

	tmp = buffer.Data[:2]
	_, err = io.ReadFull(conn, tmp)
	if err != nil {
		return err
	}
	req.Port = binary.BigEndian.Uint16(tmp[:2])
	return nil
}

func (req *Socks5Request) getAddr(conn net.Conn) error {
	buffer := buffer.SmallBytesPool.Get().(*buffer.SmallBytes)
	defer buffer.Release()

	switch req.Atyp {
	case AtypIPV4:
		tmp := buffer.Data[:4]
		_, err := io.ReadFull(conn, tmp)
		if err != nil {
			return err
		}
		req.Addr = fmt.Sprintf("%d.%d.%d.%d",
			tmp[0],
			tmp[1],
			tmp[2],
			tmp[3])
	case AtypDomain:
		tmp := buffer.Data[:1]
		_, err := io.ReadFull(conn, tmp)
		if err != nil {
			return err
		}

		len := int(tmp[0])
		tmp = make([]byte, len)
		_, err = io.ReadFull(conn, tmp)
		if err != nil {
			return err
		}
		req.Addr = string(tmp[:len])
	default:
		return ErrMessageInvalid
	}

	return nil
}

type Socks5Reply struct {
	//	+----+-----+-------+------+----------+----------+
	//  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	//  +----+-----+-------+------+----------+----------+
	//  | 1  |  1  | X'00' |  1   | Variable |    2     |
	//  +----+-----+-------+------+----------+----------+
	reply []byte
}

func NewSocksReply() *Socks5Reply {
	return &Socks5Reply{
		reply: []byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0},
	}
}

func (r *Socks5Reply) Bytes() []byte {
	return r.reply
}
