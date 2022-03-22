package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

type TlsConfigMaker struct{}

func NewTlsConfigMaker() *TlsConfigMaker {
	return &TlsConfigMaker{}
}

func (m *TlsConfigMaker) MakeTls2Config() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair("tls_cer/server.crt", "tls_cer/server.key")
	if err != nil {
		return nil, err
	}

	f, _ := os.OpenFile("./static/key.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		KeyLogWriter: f,
	}
	return cfg, nil
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

	keyLog, err := os.OpenFile("data/log_file.key", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
		KeyLogWriter: keyLog,
	}
}
