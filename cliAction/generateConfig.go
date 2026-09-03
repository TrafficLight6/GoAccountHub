package action

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/urfave/cli/v2"
)

func GenerateConfigAction() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "Generate Config File",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.json",
				Usage:   "Config File",
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := config.Config{}
			config.SaveConfig(c.String("config"), appConfig)
			fmt.Println("✅ Config File is Generated")
			return nil
		},
	}
}
