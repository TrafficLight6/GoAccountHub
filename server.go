package main

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/router"
)

func StartServer(configPath string) config.Config {
	AppConfig := config.GetConfig(configPath)
	if AppConfig.Port == "" {
		return config.Config{}
	}
	router, db := router.ReturnRouter(AppConfig)
	if router == nil || db == nil {
		return config.Config{}
	}

	//Defer Close DB Connection
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("Error when Close DB:", err)
			return
		}
		sqlDB.Close()
	}()

	router.Run(":" + AppConfig.Port)
	return AppConfig
}
