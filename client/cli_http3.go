package client

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/glstr/gwatcher/util"
	"github.com/lucas-clemente/quic-go/http3"
)

type Http3Client struct {
	addr string
}

func NewHttp3Client(addr string) *Http3Client {
	return &Http3Client{
		addr: addr,
	}
}

func (c *Http3Client) Start() error {
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			//RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       nil,
		},
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	addr := "127.0.0.1:8888/helloworld"

	rsp, err := hclient.Get(addr)
	if err != nil {
		//log.Fatal(err)
		util.Notice("get failed, error_msg:%s", err.Error())
		return err
	}
	log.Printf("Got response for %s: %#v", addr, rsp)
	defer rsp.Body.Close()
	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		util.Notice("get body failed, error_msg:%s", err.Error())
		return err
	}
	util.Notice("Request Body:%s", body.Bytes())
	return nil
}
