package middleware

import (
	"encoding/json"
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminPermissionCheckMiddleware(fieldName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get DB Handle
		db := c.Value("db").(*gorm.DB)
		//Get Cookie
		adminToken, err := c.Cookie("admin_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Cookie admin_token Not Found"})
			c.Abort()
			return
		}
		//Get Admin from Token
		admin, err := sqlOperator.GetAdminByToken(db, adminToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Admin Info Not Found"})
			c.Abort()
			return
		}
		//Check is Root Admin
		if admin.Username == "root" && admin.ID == 0 {
			c.Next()
			return
		}
		//Check Normal Admin Permission
		b, err := json.Marshal(admin.Permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Permission marshal failed"})
			c.Abort()
			return
		}
		var permission map[string]bool
		if err := json.Unmarshal(b, &permission); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Permission unmarshal failed"})
			c.Abort()
			return
		}
		//Is Permission Field Exist
		if _, ok := permission[fieldName]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "error": "Permission Denied, Because Not Exist Field: " + fieldName})
			c.Abort()
			return
		}
		//Check Permission Field
		if !permission[fieldName] {
			c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "error": "Permission Denied, Because No Permission: " + fieldName})
			c.Abort()
			return
		}
		c.Next()
	}
}
