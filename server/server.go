package server

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/router"
)

func StartServer(config config.Config) error {
	if config.Port == "" {
		return fmt.Errorf("Port is empty")
	}
	router, db := router.ReturnRouter(config)
	if router == nil || db == nil {
		return fmt.Errorf("Router or DB is nil")
	}

	//Defer Close DB Connection
	defer func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		return nil
	}()

	router.Run(":" + config.Port)
	return nil
}
