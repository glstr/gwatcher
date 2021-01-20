package server

import (
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

type QuicConnWrapper struct {
	session quic.Session
}

func NewQuicConnWrapper(session quic.Session) *QuicConnWrapper {
	w := &QuicConnWrapper{
		session: session,
	}

	//go w.run()
	return w
}

func (ss *QuicConnWrapper) run() {
	return
}

func (ss *QuicConnWrapper) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (ss *QuicConnWrapper) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (ss *QuicConnWrapper) Close() error {
	return nil
}

func (ss *QuicConnWrapper) LocalAddr() net.Addr {
	return ss.session.LocalAddr()
}

func (ss *QuicConnWrapper) RemoteAddr() net.Addr {
	return ss.session.RemoteAddr()
}

func (ss *QuicConnWrapper) SetDeadline(t time.Time) error {
	return nil
}

func (ss *QuicConnWrapper) SetReadDeadline(t time.Time) error {
	return nil
}

func (ss *QuicConnWrapper) SetWriteDeadline(t time.Time) error {
	return nil
}
