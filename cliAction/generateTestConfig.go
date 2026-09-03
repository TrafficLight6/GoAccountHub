package action

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/urfave/cli/v2"
)

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
