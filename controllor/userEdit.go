package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserEditRequestBody struct {
	UUHash string `json:"uu_hash"`

	Username string `json:"username"`
	Password string `json:"password"`
	MetaData string `json:"meta_data"`
}

func UserEdit(c *gin.Context) {
	var body UserEditRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		c.Abort()
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check Username or Password is Empty
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Username or Password Is Empty"})
		c.Abort()
		return
	}
	//Check User Exist
	var userOld sqlTable.User
	if err := db.Where("uu_hash = ?", body.UUHash).First(&userOld).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "UUHash Not Exist"})
		c.Abort()
		return
	}
	//Edit User
	userNew := userOld
	if body.Username != "" {
		userNew.Username = body.Username
	}
	if body.Password != "" {
		userNew.PasswordHash = hash.SHA256(body.Password)
	}
	if body.MetaData != "" {
		userNew.MetaData = body.MetaData
	}
	//Update User
	if err := db.Save(&userNew).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed to Save User"})
		c.Abort()
		return
	}
	//Update Same Name Character
	if userNew.Username != userOld.Username {
		var characterOld sqlTable.Character
		var characterNew sqlTable.Character
		if err := db.Where("uu_hash = ? AND character_name = ?", body.UUHash, userOld.Username).First(&characterOld).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Character Not Exist"})
			c.Abort()
			return
		}
		characterNew = characterOld
		characterNew.CharacterName = userNew.Username
		characterNew.MetaData = userNew.MetaData
		//Update Character
		if err := db.Save(&characterNew).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed to Save Character"})
			c.Abort()
			return
		}
	}
	//Return Success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success"})
	return
}
