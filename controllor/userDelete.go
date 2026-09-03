package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDeleteRequestBody struct {
	UUHash string `json:"uu_hash"`
}

func UserDelete(c *gin.Context) {
	//Get Args
	var body UserDeleteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check User Exist
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", body.UUHash).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "User Not Exist"})
		c.Abort()
		return
	}
	//Delete Character which Belong to User
	var characters []sqlTable.Character
	db.Delete(&characters, "user_uu_hash = ?", body.UUHash)
	//Delete User
	db.Delete(&user)
	//Return Success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "User Deleted"})
	return
}
