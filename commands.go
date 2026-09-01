package main

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
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
			config := StartServer(c.String("config"))
			if config.Port == "" {
				fmt.Println("⚠️ Config Error: File [" + c.String("config") + "] is Empty or Not Exist")
				return nil
			}
			return nil
		},
	}
}

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
			config.SaveConfig(c.String("config"), appConfig)
			fmt.Println("✅ Password is Updated")
			return nil
		},
	}
}
