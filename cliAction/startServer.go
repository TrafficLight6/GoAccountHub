package action

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/file"
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
			config := config.GetConfig(c.String("config"))
			//Update Frontend Port
			file.Write("./GAHFrontend/.env", "PORT="+config.FrontendPort)
			err := server.StartServer(config)
			if err != nil {
				fmt.Println("⚠️ Error when Start Server:", err)
				return err
			}
			return nil
		},
	}
}
