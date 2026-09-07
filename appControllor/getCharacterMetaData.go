package appControllor

import (
	"net/http"

	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetCharacterMetaDataRequestBody struct {
	Token string `json:"token"`
}

func GetCharacterMetaData(c *gin.Context) {
	var body GetCharacterMetaDataRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "invalid request"})
		return
	}
	//Get db
	db := c.Value("db").(*gorm.DB)
	character, err := sqlOperator.GetCharacterByToken(db, body.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "meta_data": character.MetaData})
}
