package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CharacterEditRequestBody struct {
	UserUUHash      string `json:"user_uu_hash"`
	CharacterUUHash string `json:"character_uu_hash"`

	CharacterName string `json:"character_name"`
	Password      string `json:"password"`
	MetaData      string `json:"meta_data"`
}

func CharacterEdit(c *gin.Context) {
	var body CharacterEditRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"code": http.StatusBadRequest, "error": "request body is invalid"})
		c.Abort()
		return
	}
	//Check CharacterName and Password are not empty
	if body.CharacterName == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Name or Password Is Empty"})
		c.Abort()
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check User UUHash
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", body.UserUUHash).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "UserUUHash Is Invalid"})
		c.Abort()
		return
	}
	//Get Target Character
	var characterOld sqlTable.Character
	if err := db.Where("uu_hash = ?", body.CharacterUUHash).First(&characterOld).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "CharacterUUHash Is Invalid"})
		c.Abort()
		return
	}
	//is Name Already Used By Other Character
	var characterExist sqlTable.Character
	if err := db.Where("character_name = ? AND uu_hash <> ?", body.CharacterName, body.CharacterUUHash).First(&characterExist).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Name Already Used"})
		c.Abort()
		return
	}
	//Update Character
	characterNew := characterOld
	if body.CharacterName != "" {
		characterNew.CharacterName = body.CharacterName
	}
	if body.Password != "" {
		characterNew.PasswordHash = hash.SHA256(body.Password)
	}
	if body.MetaData != "" {
		characterNew.MetaData = body.MetaData
	}
	if err := db.Save(&characterNew).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Update Character Failed"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Edit Character Success"})
}
