package adminControllor

import (
	"errors"
	"net/http"
	"time"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
	sqlOperator "github.com/TrafficLight6/GoAccountHub/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminLoginRequestBody struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsRemember bool   `json:"is_remember"`
}

func AdminLogin(c *gin.Context) {
	var body AdminLoginRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid body"})
		return
	}
	config := c.MustGet("config").(config.Config)
	_ = config
	db := c.MustGet("db").(*gorm.DB)
	_ = db
	//Root Admin Login
	if "root" == body.Username && config.RootAdminPasswordHash == hash.SHA256(body.Password) {
		token := hash.SHA256(body.Username + body.Password + time.Now().String())
		var err error
		if body.IsRemember {
			err = sqlOperator.AddAdminToken(db, config.RootAdminUUHash, token, time.Now().Add(30*24*time.Hour))
		} else {
			err = sqlOperator.AddAdminToken(db, config.RootAdminUUHash, token, time.Now().Add(1*time.Hour))
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "insert admin token failed"})
			return
		}
		//Set cookie (must be before c.JSON, otherwise headers are already flushed)
		if body.IsRemember {
			c.SetCookie("admin_token", token, 30*24*60*60, "/", "", false, true)
		} else {
			c.SetCookie("admin_token", token, 0, "/", "", false, true)
		}
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "login success", "token": token})
		return
	}
	//Normal Admin Login
	admin, err := sqlOperator.GetAdminByUsernameAndPassword(db, body.Username, hash.SHA256(body.Password))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "get admin failed"})
		return
	}
	if admin.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid credentials"})
		return
	}
	token := hash.SHA256(body.Username + body.Password + time.Now().String())
	if body.IsRemember {
		err = sqlOperator.AddAdminToken(db, admin.UUHash, token, time.Now().Add(30*24*time.Hour))
	} else {
		err = sqlOperator.AddAdminToken(db, admin.UUHash, token, time.Now().Add(1*time.Hour))
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "insert admin token failed"})
		return
	}
	//Set cookie (must be before c.JSON, otherwise headers are already flushed)
	if body.IsRemember {
		c.SetCookie("admin_token", token, 30*24*60*60, "/", "", false, true)
	} else {
		c.SetCookie("admin_token", token, 0, "/", "", false, true)
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "login success", "token": token})
}
