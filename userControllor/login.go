package userControllor

import (
	"errors"
	"net/http"
	"time"

	"github.com/TrafficLight6/GoAccountHub/hash"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	sqlTable "github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserLoginRequestBody struct {
	Username      string `json:"username"`
	CharacterName string `json:"character_name"`
	Password      string `json:"password"`
	IsRemember    bool   `json:"is_remember"`
}

func Login(c *gin.Context) {
	var body UserLoginRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid body"})
		return
	}
	//Check request body
	if body.Username == "" || body.CharacterName == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Username, character name or password is empty"})
		return
	}
	//Get db
	db := c.MustGet("db").(*gorm.DB)
	_ = db
	//Is user exist
	user, err := sqlOperator.GetUserByUsernameAndPassword(db, body.Username, hash.SHA256(body.Password))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Get user failed"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid credentials"})
		return
	}
	//Get user UUhash
	userUUHash, err := sqlOperator.GetUserUUHash(db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Get user UUhash failed"})
		return
	}
	//Get characters by user UUhash
	characters, err := sqlOperator.GetCharacterByUserUUHash(db, userUUHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Get characters failed"})
		return
	}
	//Check character name exist
	isCharacterExist := false
	characterUUHash := ""
	for _, character := range characters {
		if character.CharacterName == body.CharacterName {
			isCharacterExist = true
			characterUUHash = character.UUHash
			break
		}
	}
	if !isCharacterExist {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Character name not exist"})
		return
	}
	//Generate token
	token := hash.SHA256(body.Username + body.CharacterName + body.Password + time.Now().String())
	//Add token to db
	var userToken sqlTable.UserLoginToken
	if body.IsRemember {
		userToken = sqlTable.UserLoginToken{
			UUHash:          userUUHash,
			CharacterUUHash: characterUUHash,
			Token:           token,
			DeadTime:        time.Now().Add(30 * 24 * time.Hour),
		}
	} else {
		userToken = sqlTable.UserLoginToken{
			UUHash:          userUUHash,
			CharacterUUHash: characterUUHash,
			Token:           token,
			DeadTime:        time.Now().Add(1 * time.Hour),
		}
	}
	//Add token to db
	err = db.Create(&userToken).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Add token failed"})
		return
	}
	//Set cookie (must be before c.JSON, otherwise headers are already flushed)
	if body.IsRemember {
		c.SetCookie("user_token", token, 30*24*60*60, "/", "", false, true)
	} else {
		c.SetCookie("user_token", token, 0, "/", "", false, true)
	}
	//Return token
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Login success", "token": token})
	return
}
