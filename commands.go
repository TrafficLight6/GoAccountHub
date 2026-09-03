package main

import (
	"fmt"
	"strconv"
	"time"

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
			appConfig.RootAdminUUHash = hash.SHA256("root" + appConfig.RootAdminPasswordHash + strconv.Itoa(int(time.Now().Unix())))
			config.SaveConfig(c.String("config"), appConfig)
			fmt.Println("✅ Password is Updated")
			return nil
		},
	}
}

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
func GenerateTestConfigAction() *cli.Command {
	return &cli.Command{
		Name:  "generate-test",
		Usage: "Generate Test Config File",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.json",
				Usage:   "Config File",
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := config.Config{
				Port:                  "8080",
				DatabaseName:          "account_hub",
				DatabaseHost:          "127.0.0.1",
				DatabasePort:          "3306",
				DatabaseUser:          "postgres",
				DatabasePassword:      "postgres",
				RootAdminPasswordHash: hash.SHA256("123"),
				RootAdminUUHash:       hash.SHA256("root" + hash.SHA256("123") + strconv.Itoa(int(time.Now().Unix()))),
				SwitchConfig: config.SwitchConfig{
					AllowMultiCharacter: false,
				},
			}
			config.SaveConfig(c.String("config"), appConfig)
			fmt.Println("✅ Config File is Generated, Default Password is 123")
			return nil
		},
	}
}
