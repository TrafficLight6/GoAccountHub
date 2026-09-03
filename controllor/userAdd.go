package controllor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TrafficLight6/GoAccountHub/hash"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAddRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	MetaData string `json:"meta_data"`
}

func UserAdd(c *gin.Context) {
	var body AddUserRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//is Name Exist
	var user sqlTable.User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Username Already Exist"})
		c.Abort()
		return
	}
	//Check Username or Password is Empty
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Username or Password Is Empty"})
		c.Abort()
		return
	}
	//Add User
	user = sqlTable.User{
		Username:     body.Username,
		PasswordHash: hash.SHA256(body.Password),
		UUHash:       hash.SHA256(body.Username + hash.SHA256(body.Password) + strconv.Itoa(int(time.Now().Unix()))),
		MetaData:     body.MetaData,
	}
	err := db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal Server Error When Adding User"})
		c.Abort()
		return
	}
	//Add Same Name Character
	err = sqlOperator.AddCharacter(db, sqlTable.Character{
		CharacterName: user.Username,
		PasswordHash:  user.PasswordHash,
		UserUUHash:    user.UUHash,
		UUHash:        user.UUHash,
		MetaData:      user.MetaData,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal Server Error When Adding Character"})
		c.Abort()
		return
	}
	//Return Success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "User Added Successfully,And Same Name Character Added Successfully"})
}
