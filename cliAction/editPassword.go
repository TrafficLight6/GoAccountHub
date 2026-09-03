package action

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/urfave/cli/v2"
)

func EditPasswordAction() *cli.Command {
	return &cli.Command{
		Name:  "password",
		Usage: "Edit Root Admin Password",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.json",
				Usage:   "Config File",
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := config.GetConfig(c.String("config"))
			if c.Args().Len() == 0 {
				fmt.Println("⚠️ Password is Empty")
				return nil
			}
			appConfig.RootAdminPasswordHash = hash.SHA256(c.Args().First())
			appConfig.RootAdminUUHash = hash.SHA256("root" + appConfig.RootAdminPasswordHash + strconv.Itoa(int(time.Now().Unix())))
			config.SaveConfig(c.String("config"), appConfig)
			fmt.Println("✅ Password is Updated")
			return nil
		},
	}
}
