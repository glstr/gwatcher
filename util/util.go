package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal() {
	ch := make(chan os.Signal)
	//监听所有信号
	signal.Notify(ch, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	s := <-ch
	fmt.Println("退出信号", s)
}

// Setup a bare-bones TLS config for the server
func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	//tlsCert, err := tls.LoadX509KeyPair("conf/cert.pem", "conf/key.pem")
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
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
