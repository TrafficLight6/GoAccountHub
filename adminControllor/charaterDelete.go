package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteCharacterBody struct {
	UserUUHash      string `json:"user_uu_hash"`
	CharacterUUHash string `json:"character_uu_hash"`
}

func CharacterDelete(c *gin.Context) {
	//Get Request Body
	var body DeleteCharacterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "request body is invalid"})
		c.Abort()
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check User Exist
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", body.UserUUHash).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "User Not Exist"})
		c.Abort()
		return
	}
	//Check Character Exist
	var character sqlTable.Character
	if err := db.Where("uu_hash = ? AND user_uu_hash = ?", body.CharacterUUHash, body.UserUUHash).First(&character).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Not Exist"})
		c.Abort()
		return
	}
	//Cannot Delete Same Name Character
	if body.CharacterUUHash == body.UserUUHash {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Cannot Delete Same Name Character. If You Want To Do That,Please Turn To User Delete Api"})
		c.Abort()
		return
	}
	//Delete Character
	db.Delete(&character)
	//Return Success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Character Deleted"})
	return
}
