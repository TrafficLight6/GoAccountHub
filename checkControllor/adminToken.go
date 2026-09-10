package checkControllor

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckAdminToken(c *gin.Context) {
	cookie, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "admin_token not found",
		})
		c.Abort()
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	if !sqlOperator.CheckAdminToken(db, cookie) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "admin_token invalid",
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "admin_token valid",
	})
	return
}
