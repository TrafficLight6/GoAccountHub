package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Info(c *gin.Context) {
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Get Config
	config := c.Value("config").(config.Config)

	var (
		userCount       int64
		characterCount  int64
		userTokenCount  int64
		adminTokenCount int64
		adminCount      int64
	)

	//Count Users (soft-deleted rows are excluded automatically)
	if err := db.Model(&sqlTable.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Users"})
		return
	}
	//Count Characters
	if err := db.Model(&sqlTable.Character{}).Count(&characterCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Characters"})
		return
	}
	//Count Tokens (user tokens + admin tokens)
	if err := db.Model(&sqlTable.UserLoginToken{}).Count(&userTokenCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count User Tokens"})
		return
	}
	if err := db.Model(&sqlTable.AdminLoginToken{}).Count(&adminTokenCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Admin Tokens"})
		return
	}
	//Count Admins (root admin is not stored in table, add it back if configured)
	if err := db.Model(&sqlTable.Admin{}).Count(&adminCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Admins"})
		return
	}
	if config.RootAdminUUHash != "" {
		adminCount++
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"user_count":        userCount,
			"character_count":   characterCount,
			"total_token_count": userTokenCount + adminTokenCount,
			"user_token_count":  userTokenCount,
			"admin_token_count": adminTokenCount,
			"admin_count":       adminCount, //Not include root admin
		},
	})
}
