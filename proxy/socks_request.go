package proxy

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/glstr/gwatcher/buffer"
)

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

func (req *Socks5Request) Load(conn net.Conn) error {
	buf := buffer.GetPacketBuffer()
	defer buf.Release()
	_, err := io.ReadFull(conn, buf.Buffer[:4])
	if err != nil {
		return err
	}

	req.Ver, req.Cmd, req.Rsv, req.Atyp = buf.Buffer[0], buf.Buffer[1], buf.Buffer[2], buf.Buffer[3]
	if req.Ver != Socks5Version {
		return ErrVersionFailed
	}

	err = req.getAddr(conn)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(conn, buf.Buffer[:2])
	if err != nil {
		return err
	}
	req.Port = binary.BigEndian.Uint16(buf.Buffer[:2])
	return nil
}

func (req *Socks5Request) getAddr(conn net.Conn) error {
	buf := buffer.GetPacketBuffer()
	defer buf.Release()
	switch req.Atyp {
	case AtypIPV4:
		_, err := io.ReadFull(conn, buf.Buffer[:4])
		if err != nil {
			return err
		}
		req.Addr = fmt.Sprintf("%d.%d.%d.%d",
			buf.Buffer[0],
			buf.Buffer[1],
			buf.Buffer[2],
			buf.Buffer[3])
	case AtypDomain:
		_, err := io.ReadFull(conn, buf.Buffer[:1])
		if err != nil {
			return err
		}
		len := int(buf.Buffer[0])
		_, err = io.ReadFull(conn, buf.Buffer[:len])
		if err != nil {
			return err
		}
		req.Addr = string(buf.Buffer[:len])
	default:
		return ErrMessageInvalid
	}

	return nil
}
