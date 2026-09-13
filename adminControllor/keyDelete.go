package adminControllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KeyDeleteRequestBody struct {
	KeyUserName string `json:"key_user_name"`
}

func KeyDelete(c *gin.Context) {
	var body KeyDeleteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid request body"})
		return
	}
	//Check key user name is empty
	if body.KeyUserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "key_user_name is empty"})
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	//Delete key user name
	if err := db.Delete(&sqlTable.ApplicationKey{}, "key_user = ?", body.KeyUserName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	//Return success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success"})
}
