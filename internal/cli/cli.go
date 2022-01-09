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
	defer log.Stop()
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
		Commands: []*cli.Command{
			{
				Name: "upgrade",
				Action: func(c *cli.Context) error {
					checkUpdateFromCli()
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			var configFile = "config.yaml"
			if c.Path("config") != "" {
				configFile = c.Path("config")
			}
			conf := &server.Config{}
			tpl := &server.Config{}
			_ = yaml.Unmarshal(configTemplate, tpl)
			if content, err := os.ReadFile(configFile); err == nil {
				if err := yaml.Unmarshal(content, conf); err != nil {
					return err
				}
				if conf.Version != tpl.Version {
					backup := fmt.Sprintf("%s_%.1f", configFile, conf.Version)
					log.Warn("Old version of configuration file detected, new configuration file being generated, old configuration backed up to %s", backup)
					if err := os.WriteFile(backup, content, 0644); err != nil {
						return err
					}
					return os.WriteFile(configFile, configTemplate, 0644)
				}

			} else if os.IsNotExist(err) && configFile == "config.yaml" {
				log.Warn("Generate default configurations to config.yaml, please configure and run again.")
				return os.WriteFile("config.yaml", configTemplate, 0644)
			} else {
				return err
			}
			if c.String("addr") != "" {
				conf.Addr = c.String("addr")
			}
			if c.String("token") != "" {
				conf.Token = c.String("token")
			}
			if c.String("db") != "" {
				conf.Database = c.String("db")
			}
			if c.String("log") != "" {
				conf.LogLevel = c.String("log")
			}

			if conf.CheckUpgrade {
				checkUpdateFromCli()
			}

			server.New(conf).Run()
			return nil
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	fmt.Printf(banner, server.VERSION)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
