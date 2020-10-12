package action

import (
	"github.com/glstr/gwatcher/server"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

func StartServer(c *cli.Context) error {
	defaultAddr := "0.0.0.0:8888"
	ser := server.NewUdpServer(defaultAddr)
	go ser.Start()
	util.WaitSignal()
	return nil
}
