package action

import (
	"github.com/glstr/gwatcher/client"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

func StartClient(c *cli.Context) error {
	//var addr string = "localhost:8888"
	addr := c.String("addr")
	protocol := c.String("protocol")
	cli := client.NewClient(protocol, addr)
	go cli.Start()
	util.WaitSignal()
	return nil
}
