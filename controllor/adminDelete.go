package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminDeleteRequestBody struct {
	UUHash string `json:"uu_hash"`
}

func AdminDelete(c *gin.Context) {
	//Get Request Body
	var body AdminDeleteRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Request Body"})
		return
	}
	//Get Cookie admin_token
	cookie, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Cookie"})
		return
	}
	//Get DB
	db := c.Value("db").(*gorm.DB)
	//Get Config
	config := c.Value("config").(config.Config)
	//Check admin_token
	var adminToken sqlTable.AdminLoginToken
	err = db.Where("token = ?", cookie).First(&adminToken).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Admin Token"})
		return
	}
	//Cannot Delete Root Admin
	if body.UUHash == config.RootAdminUUHash {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Root Admin Cannot be Deleted"})
		return
	}
	//Check is not self delete
	if adminToken.UUHash == body.UUHash {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Cannot Delete Yourself"})
		return
	}
	//Delete Target Admin
	result := db.Where("uu_hash = ?", body.UUHash).Delete(&sqlTable.Admin{})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Delete Admin Failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Target Admin Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Delete Admin Success"})
}
