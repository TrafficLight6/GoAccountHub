package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CharacterGetRequestBody struct {
	UserUUHash      string `json:"user_uu_hash"`
	CharacterUUHash string `json:"character_uu_hash"`
}

func CharacterGet(c *gin.Context) {
	var body CharacterGetRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check Character Exist
	var character sqlTable.Character
	if err := db.Where("uu_hash = ? AND user_uu_hash = ?", body.CharacterUUHash, body.UserUUHash).First(&character).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Not Exist"})
		c.Abort()
		return
	}
	//Return Character
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "character": character})
}
