package action

import (
	"github.com/glstr/gwatcher/server"
	"github.com/glstr/gwatcher/util"
	"github.com/urfave/cli"
)

func StartServer(c *cli.Context) error {
	addr := "0.0.0.0:8888"
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
