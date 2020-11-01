package main

import (
	"log"
	"os"
	"sort"

	"github.com/glstr/gwatcher/action"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "gwatcher",
		Usage: "watcher who always watchs the world",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang,l",
				Value: "english",
				Usage: "Language for the greeting",
			},
			&cli.StringFlag{
				Name:  "config,c",
				Usage: "Load configuration from `FILE`",
			},
		},

		Commands: []cli.Command{
			{
				Name:    "server",
				Aliases: []string{"ser"},
				Usage:   "start a server",
				Action:  action.StartServer,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "protocol,p",
						Value: "udp",
						Usage: "server protocol, support udp & quic",
					},

					&cli.StringFlag{
						Name:  "address,addr",
						Value: "127.0.0.1:443",
						Usage: "server protocol, support udp & quic",
					},
				},
			},
			{
				Name:    "client",
				Aliases: []string{"cli"},
				Usage:   "start a client",
				Action:  action.StartClient,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "protocol, p",
						Value: "udp",
						Usage: "server protocol, support udp & quic",
					},

					&cli.StringFlag{
						Name:  "address, addr",
						Value: "127.0.0.1:443",
						Usage: "server protocol, support udp & quic",
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
