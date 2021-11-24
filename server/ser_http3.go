package server

import (
	"net/http"

	"github.com/glstr/gwatcher/util"
	"github.com/lucas-clemente/quic-go/http3"
)

type Http3Server struct {
	addr string
}

func NewHttp3Server(addr string) *Http3Server {
	return &Http3Server{
		addr: addr,
	}
}

func (s *Http3Server) Start() error {
	http.HandleFunc("/helloworld", helloHandle)

	cerPath := "./conf/cert/glstr.cer"
	keyPath := "./conf/cert/server.key"
	http3.ListenAndServeQUIC(s.addr, cerPath, keyPath, nil)
	return nil
}

func (s *Http3Server) Stop() {
}

func helloHandle(res http.ResponseWriter, req *http.Request) {
	util.Notice("receive request")
	result := `{"hello":"world"}`
	res.Write([]byte(result))
}
