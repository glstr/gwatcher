package util

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal() {
	ch := make(chan os.Signal, 1)
	//监听所有信号
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	s := <-ch
	fmt.Println("退出信号", s)
}

// A wrapper for io.Writer that also logs the message.
type LoggingWriter struct{ io.Writer }

func (w LoggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

func DisplaySocketOption(conn net.Conn) error {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		rawConn, err := tcpConn.SyscallConn()
		if err != nil {
			Notice("get raw conn failed:%s", err.Error())
			return err
		}

		tcpConn.SetWriteBuffer(5000)
		tcpConn.SetReadBuffer(1000)

		f := func(fd uintptr) {
			rdbuf, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
			if err != nil {
				Notice("get fd option failed, error_msg:%s", err.Error())
				return
			}

			sdbuf, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
			if err != nil {
				Notice("get fd option failed, error_msg:%s", err.Error())
				return
			}
			Notice("rbuf:%d, sdbuf:%d", rdbuf, sdbuf)
		}

		return rawConn.Control(f)
	}
	return nil
}
