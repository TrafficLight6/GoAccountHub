package router

import (
	"github.com/TrafficLight6/GoAccountHub/controllor"
	"github.com/gin-gonic/gin"
)

func ReturnRouter() *gin.Engine {
	router := gin.Default()
	router.RouterGroup.Group("/")
	//Root Page
	router.GET("/", controllor.Root)

	return router
}
