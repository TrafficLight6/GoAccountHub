package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CharacterCount(c *gin.Context) {
	//Get Db
	db := c.Value("db").(*gorm.DB)

	//Count Characters (soft-deleted rows are excluded automatically)
	var characterCount int64
	if err := db.Model(&sqlTable.Character{}).Count(&characterCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed To Count Characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"character_count": characterCount,
		},
	})
}
