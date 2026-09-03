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
		v1.POST("/admin/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_add_admin"), controllor.AdminAdd)
		v1.DELETE("/admin/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_delete_admin"), controllor.AdminDelete)
		v1.PUT("/admin/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_edit_admin"), controllor.AdminEdit)
		v1.GET("/admin/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_get_admin"), controllor.AdminGet)

		//User Api
		//User Operation, Admin Only
		v1.POST("/user/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), controllor.UserAdd)
		v1.DELETE("/user/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), controllor.UserDelete)
		v1.PUT("/user/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), controllor.UserEdit)
		v1.GET("/user/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), controllor.UserGet)

		//Character Api
		//Also Admin Only
		//Check is Allowed to Use Multi Character Functions
		v1.Use(middleware.ConfigBlocker("allow_multi_character"))
		//Character Operation
		v1.POST("/character/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), controllor.CharacterAdd)
		v1.PUT("/character/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), controllor.CharacterEdit)
	}

	return router, db
}
