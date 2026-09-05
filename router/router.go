package router

import (
	"github.com/TrafficLight6/GoAccountHub/adminControllor"
	"github.com/TrafficLight6/GoAccountHub/appControllor"
	"github.com/TrafficLight6/GoAccountHub/config"
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
	router.GET("/", adminControllor.Root)

	//v1 Api
	v1 := router.Group("/api/v1")
	{
		//ALL ADMIN ONLY

		//Admin Api
		//Login & Logout
		v1.POST("/admin/login", adminControllor.AdminLogin)
		v1.DELETE("/admin/logout", adminControllor.AdminLogout)

		//Admin Operation
		v1.POST("/admin/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_add_admin"), adminControllor.AdminAdd)
		v1.DELETE("/admin/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_delete_admin"), adminControllor.AdminDelete)
		v1.PUT("/admin/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_edit_admin"), adminControllor.AdminEdit)
		v1.GET("/admin/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_get_admin"), adminControllor.AdminGet)
		v1.GET("/admin/range", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_get_admin_range"), adminControllor.AdminRange)

		//User Api
		v1.POST("/user/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), adminControllor.UserAdd)
		v1.DELETE("/user/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), adminControllor.UserDelete)
		v1.PUT("/user/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), adminControllor.UserEdit)
		v1.GET("/user/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user"), adminControllor.UserGet)
		v1.GET("/user/range", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_user_range"), adminControllor.UserRange)

		//Character Api
		//Check is Allowed to Use Multi Character Functions
		v1.Use(middleware.ConfigBlocker("allow_multi_character"))
		//Character Operation
		v1.POST("/character/add", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), adminControllor.CharacterAdd)
		v1.DELETE("/character/delete", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), adminControllor.CharacterDelete)
		v1.PUT("/character/edit", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), adminControllor.CharacterEdit)
		v1.GET("/character/get", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character"), adminControllor.CharacterGet)
		v1.GET("/character/range", middleware.AdminCheckMiddleware(), middleware.AdminPermissionCheckMiddleware("can_operate_character_range"), adminControllor.CharacterRange)
	}

	{
		//USER ONLY
		export := v1.Group("/app")
		//Login & Logout
		export.POST("/user/login", middleware.ApplicationKeyCheckMiddleware(), appControllor.Login)
		export.DELETE("/user/logout", middleware.ApplicationKeyCheckMiddleware(), appControllor.Logout)
	}
	return router, db
}
