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
			action.ServerCmd,
			action.ClientCmd,
			action.CodecCmd,
			action.ProxyCmd,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
