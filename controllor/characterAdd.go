package controllor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CharacterAddRequestBody struct {
	CharacterName string `json:"character_name"`
	Password      string `json:"password"`
	UserUUHash    string `json:"user_uu_hash"`
	MetaData      string `json:"meta_data"`
}

type CharacterAddResponseBody struct {
	CharacterName string `json:"character_name"`
	Password      string `json:"password"`
	UserUUHash    string `json:"user_uu_hash"`
	MetaData      string `json:"meta_data"`
}

func CharacterAdd(c *gin.Context) {
	var body CharacterAddRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"code": http.StatusBadRequest, "error": "request body is invalid"})
		c.Abort()
		return
	}
	//Check Body
	if body.CharacterName == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Name, Password Is Empty"})
		c.Abort()
		return
	}
	if body.UserUUHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "UserUUHash Is Empty"})
		c.Abort()
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check UserUUHash is vaild or exist
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", body.UserUUHash).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "UserUUHash Is Invalid"})
		c.Abort()
		return
	}
	//is Name Already Used
	var character sqlTable.Character
	if err := db.Where("character_name = ?", body.CharacterName).First(&character).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Name Already Used"})
		c.Abort()
		return
	}
	//Add Character
	character = sqlTable.Character{
		CharacterName: body.CharacterName,
		PasswordHash:  hash.SHA256(body.Password),
		UserUUHash:    body.UserUUHash,
		UUHash:        hash.SHA256(body.UserUUHash + body.CharacterName + body.Password + strconv.Itoa(int(time.Now().Unix()))),
		MetaData:      body.MetaData,
	}
	err := db.Create(&character).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal Server Error When Adding Character"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Character Added Successfully", "character": character})
}
