package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// rootPermission returns all permissions for the root admin.
func rootPermission() gin.H {
	return gin.H{
		"can_add_admin":         true,
		"can_delete_admin":      true,
		"can_edit_admin":        true,
		"can_get_admin":         true,
		"can_operate_user":      true,
		"can_operate_character": true,
	}
}

func AdminInfo(c *gin.Context) {
	cookie, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Admin Token Not Found",
		})
		return
	}
	//Get db & config
	db := c.Value("db").(*gorm.DB)
	appConfig := c.Value("config").(config.Config)

	var tokenRecord sqlTable.AdminLoginToken
	if err := db.Where("token = ?", cookie).First(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Admin Token Invalid",
		})
		return
	}
	//Root admin: stored in config, not in admins table; root has all permissions
	if tokenRecord.UUHash == appConfig.RootAdminUUHash {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Admin Info",
			"data": gin.H{
				"username":   "root",
				"uu_hash":    appConfig.RootAdminUUHash,
				"is_root":    true,
				"permission": rootPermission(),
			},
		})
		return
	}
	//Normal admin
	var admin sqlTable.Admin
	if err := db.Where("uu_hash = ?", tokenRecord.UUHash).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Admin Info Not Found",
		})
		return
	}
	//Only return safe fields (never expose PasswordHash)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Admin Info",
		"data": gin.H{
			"username":   admin.Username,
			"uu_hash":    admin.UUHash,
			"permission": admin.Permission,
			"is_root":    false,
		},
	})
}
