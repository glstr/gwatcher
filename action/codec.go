package action

import (
	"github.com/glstr/gwatcher/protocol"
	"github.com/glstr/gwatcher/server"
	"github.com/urfave/cli"
)

var CodecCmd = cli.Command{
	Name:    "codec",
	Aliases: []string{"codec"},
	Usage:   "parse packet for pointed protocol",
	Action:  StartCodec,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "protocol, p",
			Value: "quic",
			Usage: "server protocol, support:" + server.DisplayProtocols(),
		},
		&cli.StringFlag{
			Name:  "file, f",
			Value: "./data.txt",
			Usage: "data file path",
		},
	},
}

func StartCodec(c *cli.Context) error {
	filePath := c.String("file")
	p := c.String("protocol")

	option := &protocol.CodecOption{
		FilePath: filePath,
		Protocol: protocol.ProtocolType(p),
	}

	codec := protocol.NewCodec()

	return codec.Parse(option)
}
