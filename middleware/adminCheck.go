package middleware

import (
	"net/http"
	"time"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("admin_token")
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "admin_token not found"})
			c.Abort()
			return
		}
		adminToken, err := sqlOperator.GetAdminToken(c.Value("db").(*gorm.DB), cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "admin_token invalid"})
			c.Abort()
			return
		}
		//Check token expiration
		if !adminToken.DeadTime.After(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "admin_token expired"})
			c.Abort()
			return
		}
		c.Set("adminToken", adminToken)
		c.Next()
	}
}
