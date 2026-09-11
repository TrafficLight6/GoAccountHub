package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserCount(c *gin.Context) {
	//Get Db
	db := c.Value("db").(*gorm.DB)

	//Count Users (soft-deleted rows are excluded automatically)
	var userCount int64
	if err := db.Model(&sqlTable.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"user_count": userCount,
		},
	})
}
