package adminControllor

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminLogout(c *gin.Context) {
	//Get cookie
	token, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Token is empty"})
		return
	}

	//Validate Token
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Token is empty"})
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	//Check Token Exist
	if !sqlOperator.CheckAdminToken(db, token) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Token not found"})
		return
	}
	//Delete Token From Database(Logout)
	if err := sqlOperator.DelAdminToken(db, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Logout success"})
}
