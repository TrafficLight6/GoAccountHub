package adminControllor

import (
	"net/http"
	"regexp"
	"time"

	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KeyAddRequestBody struct {
	KeyUserName string `json:"key_user_name"`
}

func KeyAdd(c *gin.Context) {
	var body KeyAddRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid request body"})
		return
	}
	//Check name is empty
	if body.KeyUserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Key user name is empty"})
		return
	}
	//Check name is ascii string
	if !regexp.MustCompile(`^[\x20-\x7E]+$`).MatchString(body.KeyUserName) {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Key user name must be in printable ASCII"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//Check key user name is exist
	var key sqlTable.ApplicationKey
	if err := db.Where("key_user = ?", body.KeyUserName).First(&key).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Key user name is already exist"})
		return
	}
	//Add key user name
	newKey := sqlTable.ApplicationKey{
		KeyUser: body.KeyUserName,
		Key:     hash.SHA256(body.KeyUserName + time.Now().String()),
	}
	//Save key user name
	if err := db.Create(&newKey).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Failed to add key user name"})
		return
	}
	//Return key user name
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "key_user_name": newKey.KeyUser, "key": newKey.Key})
}
