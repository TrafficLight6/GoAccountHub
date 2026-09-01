package router

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/controllor"
	"github.com/TrafficLight6/GoAccountHub/middleware"
	"github.com/gin-gonic/gin"
)

func ReturnRouter(config config.Config) *gin.Engine {
	router := gin.Default()
	router.RouterGroup.Group("/")
	//Root Page
	router.GET("/", controllor.Root)

	router.Use(middleware.ConfigInsertMiddleware(config))
	return router
}
