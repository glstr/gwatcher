package action

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/glstr/gwatcher/udpserver"
	"github.com/urfave/cli"
)

func StartServer(c *cli.Context) error {
	defaultAddr := "0.0.0.0:8888"
	server := udpserver.NewUdpServer(defaultAddr)
	go server.Start()

	ch := make(chan os.Signal)
	//监听所有信号
	signal.Notify(ch)
	//阻塞直到有信号传入
	fmt.Println("启动")
	s := <-ch
	fmt.Println("退出信号", s)
	go server.Stop()
	return nil
}
