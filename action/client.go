package action

import (
	"github.com/glstr/gwatcher/client"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

var ClientCmd = cli.Command{
	Name:    "client",
	Aliases: []string{"cli"},
	Usage:   "start a client",
	Action:  StartClient,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "protocol, p",
			Value: "udp",
			Usage: "server protocol, support udp & quic",
		},

		&cli.StringFlag{
			Name:  "address, addr",
			Value: "127.0.0.1:8888",
			Usage: "server protocol, support udp & quic",
		},
	},
}

func StartClient(c *cli.Context) error {
	//var addr string = "localhost:8888"
	addr := c.String("addr")
	protocol := c.String("protocol")
	cli := client.NewClient(protocol, addr)
	go cli.Start()
	util.WaitSignal()
	return nil
}
