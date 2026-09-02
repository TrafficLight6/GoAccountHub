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

	//Middleware
	router.Use(middleware.ConfigInsertMiddleware(config))
	router.Use(middleware.DBInsertMiddleware(db))

	router.RouterGroup.Group("/")
	//Root Page
	router.GET("/", controllor.Root)

	//v1 Api
	v1 := router.Group("/api/v1")
	{
		//Admin Api
		v1.GET("/admin/login", controllor.AdminLogin)
	}

	return router, db
}
