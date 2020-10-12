package action

import (
	"github.com/glstr/gwatcher/client"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

func StartClient(c *cli.Context) error {
	addr := "127.0.0.1:8888"
	cli := client.NewUdpClient(addr)
	go cli.Start()
	util.WaitSignal()
	return nil
}
