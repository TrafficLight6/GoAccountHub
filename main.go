package main

import (
	"github.com/TrafficLight6/GoAccountHub/router"
)

func main() {
	router := router.ReturnRouter()
	router.Run(":8080")
}
