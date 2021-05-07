package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/li4n0/revsuit/pkg/server"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	log "unknwon.dev/clog/v2"
)

func Start() {
	app := &cli.App{
		Name:  "RevSuit",
		Usage: "An Open-Sourced Reverse Platform Designed for Receive Various Kinds of Connection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Usage: "platform listened address",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "token used to manage platform",
			},
			&cli.StringFlag{
				Name:  "flags",
				Usage: "for http and dns connection, platform will only record the connection those match these regex flags. * meaning record all",
			},
			&cli.StringFlag{
				Name:  "db",
				Usage: "database file path",
			},
		},
		Action: func(c *cli.Context) error {
			conf := &server.Config{}
			if content, err := os.ReadFile("config.yaml"); err == nil {
				if err := yaml.Unmarshal(content, conf); err != nil {
					return err
				}
			} else {
				return err
			}
			if c.String("addr") != "" {
				conf.Addr = c.String("addr")
			}
			if c.String("token") != "" {
				conf.Addr = c.String("token")
			}
			if c.String("db") != "" {
				conf.Database = c.String("db")
			}
			server.New(conf).Run()
			return nil
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	fmt.Println(banner)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
