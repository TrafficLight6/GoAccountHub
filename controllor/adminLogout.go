package controllor

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//db := sqlOperator.GetDB()

type AdminLogoutRequest struct {
	Token string `json:"token"`
}

func AdminLogout(c *gin.Context) {
	var body AdminLogoutRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid request body"})
		return
	}
	//Validate Token
	if body.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Token is empty"})
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	//Check Token Exist
	if !sqlOperator.CheckAdminToken(db, body.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "error": "Token not found"})
		return
	}
	//Delete Token From Database(Logout)
	if err := sqlOperator.DelAdminToken(db, body.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Logout success"})
}
