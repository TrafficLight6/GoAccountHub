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
		//Login & Logout
		v1.GET("/admin/login", controllor.AdminLogin)
		v1.DELETE("/admin/logout", controllor.AdminLogout)
		//Admin Operation
		v1.POST("/admin/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_add_admin"), controllor.AddAdmin)
		v1.DELETE("/admin/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_delete_admin"), controllor.DeleteAdmin)
		v1.PUT("/admin/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_edit_admin"), controllor.EditAdmin)
		v1.GET("/admin/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_get_admin"), controllor.GetAdmin)
		//User Api
		//User Operation
		v1.POST("/user/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), controllor.AddUser)
	}

	return router, db
}
