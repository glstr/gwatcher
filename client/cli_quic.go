package client

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/glstr/gwatcher/msg"
	"github.com/glstr/gwatcher/util"
	quic "github.com/lucas-clemente/quic-go"
)

type QuicClient struct {
	addr string
	done chan struct{}
}

func NewQuicClient(addr string) *QuicClient {
	return &QuicClient{
		addr: addr,
		done: make(chan struct{}),
	}
}

func MakeClientQuicConf() *quic.Config {
	return &quic.Config{
		MaxIncomingStreams: 2,
	}
}

func (c *QuicClient) getSession() (quic.Session, error) {
	util.Notice("start quic client")
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	return quic.DialAddr(c.addr, tlsConf, MakeClientQuicConf())
}

func (c *QuicClient) Start() error {
	for {
		select {
		case <-c.done:
			return nil
		default:
		}
		sess, err := c.getSession()
		if err != nil {
			util.Notice("get session failed:%s", err.Error())
			return err
		}

		err = c.handleSession(sess)
		if err != nil {
			util.Notice("handle session failed:%s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func (c *QuicClient) handleSession(sess quic.Session) error {
	defer func() {
		sess.CloseWithError(quic.ApplicationErrorCode(0), "application error")
	}()

	go func() {
		for {
			select {
			case <-c.done:
				util.Notice("write loop exit")
				return
			default:
			}
			stream, err := sess.OpenStreamSync(context.Background())
			if err != nil {
				util.Notice("open stream failed, error_msg:%s", err.Error())
				return
			}

			msgContainer := &msg.MessageContainer{
				Id:   uint64(time.Now().Unix()),
				Data: strings.Repeat("h", 10*1024),
			}
			parser := msg.NewParser()
			err = parser.Marshal(stream, msgContainer)
			stream.Close()
			if err != nil {
				util.Notice("write stream failed, error_msg:%s", err.Error())
				return
			}

			util.Notice("client send len:%d", len(msgContainer.Data))
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case <-c.done:
			util.Notice("read loop exit")
			return nil
		default:
		}

		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			c.Stop()
			util.Notice("accept stream failed:%s", err.Error())
			return err
		}

		parser := msg.NewParser()
		msg, err := parser.Unmarshal(stream)
		stream.Close()
		if err != nil {
			c.Stop()
			util.Notice("read data failed:%s", err.Error())
			return err
		}

		util.Notice("client read len:%d", len(msg.Data))
	}
}

func (c *QuicClient) Stop() {
	close(c.done)
}
