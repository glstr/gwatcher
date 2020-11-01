package server

import "net/http"

type HttpServer struct {
	addr
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func (s *HttpServer) Start() error {
	http.HandleFunc("/helloworld", helloHandle)

	http.ListenAndServe(s.addr, nil)
	return nil
}

func (s *HttpServer) Stop() {
	return
}
