package main

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/router"
)

func StartServer(configPath string) config.Config {
	AppConfig := config.GetConfig(configPath)
	if AppConfig.Port == "" {
		return config.Config{}
	}
	router := router.ReturnRouter(AppConfig)
	router.Run(":" + AppConfig.Port)
	return AppConfig
}
