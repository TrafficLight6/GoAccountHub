package action

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/urfave/cli/v2"
)

func DeleteApplicationKeyAction() *cli.Command {
	return &cli.Command{
		Name:  "del-key",
		Usage: "Delete Application Key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.json",
				Usage:   "Config File",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Value:   "default",
				Usage:   "Application Key",
			},
		},
		Action: func(c *cli.Context) error {
			appConfig := config.GetConfig(c.String("config"))
			db := sqlOperator.ConnectDB(appConfig)
			defer func() {
				sqlDB, err := db.DB()
				if err != nil {
					fmt.Println("Error when Close DB:", err)
					return
				}
				sqlDB.Close()
			}()
			key := c.String("key")
			err := sqlOperator.DeleteApplicationKey(db, key)
			if err != nil {
				return err
			}
			fmt.Println("✅ Delete Key Success")
			return err
		},
	}
}
