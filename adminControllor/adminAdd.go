package adminControllor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminAddRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Permission sqlTable.Permission `json:"permission"`
}

func AdminAdd(c *gin.Context) {
	var body AdminAddRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Db
	db := c.Value("db").(*gorm.DB)
	//is Name Exist
	var admin sqlTable.Admin
	if err := db.Where("username = ?", body.Username).First(&admin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Username Already Exist"})
		c.Abort()
		return
	}
	//is being Named Root
	if body.Username == "root" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Root Admin Can Be Named"})
		c.Abort()
		return
	}
	//Check Username or Password is Empty
	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Username or Password Is Empty"})
		c.Abort()
		return
	}
	//Add Admin
	err := db.Create(&sqlTable.Admin{
		Username:     body.Username,
		PasswordHash: hash.SHA256(body.Password),
		UUHash:       hash.SHA256(body.Username + hash.SHA256(body.Password) + strconv.Itoa(int(time.Now().Unix()))),

		Permission: body.Permission,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal Server Error When Adding Admin"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Admin Added", "admin": admin})
}
