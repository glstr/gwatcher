package action

import (
	"github.com/glstr/gwatcher/proxy"
	"github.com/urfave/cli"
)

var ProxyCmd = cli.Command{
	Name:    "proxy",
	Aliases: []string{"proxy"},
	Usage:   "start a proxy in ip protocol",
	Action:  StartProxy,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
			Usage: "proxy host, default: 0.0.0.0",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "8887",
			Usage: "proxy port, default: 8887",
		},
		cli.IntFlag{
			Name:  "type, t",
			Value: 0,
			Usage: "proxy type, default: 0(socks5)",
		},
	},
}

func StartProxy(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.String("port")
	proxyType := ctx.Int("type")

	config := proxy.ProxyConfig{
		Host: host,
		Port: port,
	}

	proxy, err := proxy.NewProxy(proxy.ProxyType(proxyType), &config)
	if err != nil {
		return err
	}

	return proxy.Start()
}
