package appControllor

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetUserMetaDataRequestBody struct {
	Token string `json:"token"`
}

func GetUserMetaData(c *gin.Context) {
	var body GetUserMetaDataRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid body"})
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	//Get user by token
	user, err := sqlOperator.GetUserByToken(db, body.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success", "meta_data": user.MetaData})
}
