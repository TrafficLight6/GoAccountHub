package action

import (
	"fmt"
	"time"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/urfave/cli/v2"
)

func AddApplicationKeyAction() *cli.Command {
	return &cli.Command{
		Name:  "add-key",
		Usage: "Add Application Key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.json",
				Usage:   "Config File",
			},
			&cli.StringFlag{
				Name:    "site",
				Aliases: []string{"s"},
				Value:   "default",
				Usage:   "Application Site Site",
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
			var applicationKey sqlTable.ApplicationKey
			applicationKey.ApplicationUsingSite = c.String("site")
			applicationKey.Key = hash.SHA256(c.String("site") + time.Now().String())
			err := sqlOperator.AddApplicationKey(db, applicationKey)
			if err != nil {
				return err
			}
			fmt.Println("✅ Add Key Success")
			return err
		},
	}
}
