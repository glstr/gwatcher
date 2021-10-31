package action

import (
	"github.com/glstr/gwatcher/server"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

var ServerCmd = cli.Command{
	Name:    "server",
	Aliases: []string{"ser"},
	Usage:   "start a server",
	Action:  StartServer,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "protocol, p",
			Value: "udp",
			Usage: "server protocol, support:" + server.DisplayProtocols(),
		},
		&cli.StringFlag{
			Name:  "address, addr",
			Value: "0.0.0.0:8888",
			Usage: "server address",
		},
		&cli.StringFlag{
			Name:  "handler, hd",
			Value: string(server.HEcho),
			Usage: "server handler, support:" + server.DisplayHandlers(),
		},
		&cli.StringFlag{
			Name:  "server_type, st",
			Value: "default",
			Usage: "server type, support:",
		},
	},
}

func StartServer(c *cli.Context) error {
	addr := c.String("address")
	protocol := c.String("protocol")
	hType := c.String("handler")
	ser, err := server.GetServer(
		server.ProtocolType(protocol),
		server.HandlerType(hType),
		addr)
	if err != nil {
		util.Notice("get server failed:%s", err.Error())
		return err
	}

	go ser.Start()
	util.WaitSignal()
	return nil
}
