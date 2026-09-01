package main

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/router"
)

func StartServer(configPath string) config.Config {
	router := router.ReturnRouter()
	AppConfig := config.GetConfig(configPath)
	if AppConfig.Port == "" {
		return config.Config{}
	}
	router.Run(":" + AppConfig.Port)
	return AppConfig
}
