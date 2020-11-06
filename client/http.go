package client

import (
	"net/http"
	"time"

	"github.com/glstr/gwatcher/util"
)

type HttpClient struct {
	addr string
}

func NewHttpClient(addr string) *HttpClient {
	return &HttpClient{
		addr: addr,
	}
}

func (c *HttpClient) Start() error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	url := "http://" + c.addr + "/helloworld"
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	util.Notice("resp:%v", resp)
	return nil
}
