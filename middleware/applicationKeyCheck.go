package middleware

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApplicationKeyCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := c.Cookie("app_key")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  http.StatusBadRequest,
				"error": "app_key cookie is required",
			})
			c.Abort()
			return
		}
		db := c.Value("db").(*gorm.DB)
		if !sqlOperator.CheckApplicationKey(db, key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  http.StatusBadRequest,
				"error": "app_key cookie is invalid",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
