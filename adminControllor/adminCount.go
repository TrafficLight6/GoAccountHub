package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminCount(c *gin.Context) {
	//Get Db
	db := c.Value("db").(*gorm.DB)

	//Count Admins (soft-deleted rows are excluded automatically)
	var adminCount int64
	if err := db.Model(&sqlTable.Admin{}).Count(&adminCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Admins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"admin_count": adminCount,
		},
	})
}
