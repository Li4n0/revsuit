package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/li4n0/revsuit/pkg/server"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	log "unknwon.dev/clog/v2"
)

func Start() {
	_ = log.NewConsole(100,
		log.ConsoleConfig{
			Level: log.LevelInfo,
		})
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
				Name:  "db",
				Usage: "database file path",
			},
			&cli.StringFlag{
				Name:  "log",
				Usage: "specify a log level from debug, info, warn, error, fatal. (default: info)",
			},
			&cli.PathFlag{
				Name:  "config",
				Usage: "load configuration from FILE (default: config.yaml)",
			},
		},
		Action: func(c *cli.Context) error {
			var configFile = "config.yaml"
			if c.Path("config") != "" {
				configFile = c.Path("config")
			}
			conf := &server.Config{}
			if content, err := os.ReadFile(configFile); err == nil {
				if err := yaml.Unmarshal(content, conf); err != nil {
					return err
				}
			} else if os.IsNotExist(err) && configFile == "config.yaml" {
				log.Warn("Generate default configurations to config.yaml, please configure and run again.")
				err := os.WriteFile("config.yaml", configTemplate, 0644)
				if err != nil {
					return err
				}
				return nil
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
			if c.String("log") != "" {
				conf.Database = c.String("log")
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
