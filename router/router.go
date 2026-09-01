package router

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/controllor"
	"github.com/TrafficLight6/GoAccountHub/middleware"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReturnRouter(config config.Config) (*gin.Engine, *gorm.DB) {
	router := gin.Default()
	db := sqlOperator.ConnectDB(config)
	if db == nil {
		return nil, nil
	}

	router.RouterGroup.Group("/")
	//Root Page
	router.GET("/", controllor.Root)

	//Middleware
	router.Use(middleware.ConfigInsertMiddleware(config))
	router.Use(middleware.DBInsertMiddleware(db))
	return router, db
}
