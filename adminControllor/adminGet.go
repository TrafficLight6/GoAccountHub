package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminGetResponseBody struct {
	UUHash string `json:"uu_hash"`
}

func AdminGet(c *gin.Context) {
	//Get Request Body
	var body AdminGetResponseBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Get Config
	config := c.Value("config").(config.Config)

	//Is Root Admin
	if config.RootAdminUUHash == body.UUHash {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "admin": nil, "is_root": true})
		return
	}
	//Get Target Admin Info
	var admin sqlTable.Admin
	err := db.Where("uu_hash = ?", body.UUHash).First(&admin).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "Admin Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "admin": admin, "is_root": false})
}
