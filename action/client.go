package action

import (
	"github.com/glstr/gwatcher/client"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

func StartClient(c *cli.Context) error {
	var addr string = "127.0.0.1:8888"
	if c.String("address") != "" {
		addr = c.String("address")
	}
	cli := client.NewUdpClient(addr)
	go cli.Start()
	util.WaitSignal()
	return nil
}
