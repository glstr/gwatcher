package server

import (
	"net/http"
	"time"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func (s *HttpServer) Start() error {
	http.HandleFunc("/helloworld", helloHandle)

	server := http.Server{
		Addr:        s.addr,
		IdleTimeout: 5 * time.Second,
	}
	//http.ListenAndServe(s.addr, nil)
	server.ListenAndServe()
	return nil
}

func (s *HttpServer) Stop() {
	return
}
