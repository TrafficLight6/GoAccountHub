package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserGetRequestBody struct {
	UUhash string `json:"uu_hash"`
}

func UserGet(c *gin.Context) {
	var body UserGetRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request"})
		c.Abort()
		return
	}
	//Get DB
	db := c.Value("db").(*gorm.DB)
	//Check User Exist
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", body.UUhash).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "User Not Exist"})
		c.Abort()
		return
	}
	//Get Belonging Character Number
	var characterNumber int64
	if err := db.Model(&sqlTable.Character{}).Where("user_uu_hash = ?", user.UUHash).Count(&characterNumber).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Number Get Failed"})
		c.Abort()
		return
	}
	//Return Character Number
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success", "user": user, "character_number": characterNumber})
}
