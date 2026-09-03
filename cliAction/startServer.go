package action

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/server"
	"github.com/urfave/cli/v2"
)

func StartServerAction() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start GoAccountHub",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.json",
				Usage:   "config file",
			},
		},
		Action: func(c *cli.Context) error {
			config := server.StartServer(c.String("config"))
			if config.Port == "" {
				fmt.Println("⚠️ Config Error: File [" + c.String("config") + "] is Empty or Not Exist")
				return nil
			}
			return nil
		},
	}
}
