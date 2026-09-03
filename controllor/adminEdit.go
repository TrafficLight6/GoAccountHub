package controllor

import (
	"net/http"

	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/hash"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EditAdminRequestBody struct {
	UUHash    string                `json:"uu_hash"`    //Target Admin uu_hash
	AdminInfo sqlTable.AdminForEdit `json:"admin_info"` //Target Admin Info
}

func EditAdmin(c *gin.Context) {
	//Get Request Body
	var body EditAdminRequestBody
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
	//Cannot Edit Root Admin
	if body.UUHash == config.RootAdminUUHash {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Root Admin Cannot be Edited"})
		return
	}
	//Check is not self edit
	if adminToken.UUHash == body.UUHash {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Cannot Edit Yourself"})
		return
	}
	//Get Target Admin Old Info
	var adminOld sqlTable.Admin
	err = db.Where("uu_hash = ?", body.UUHash).First(&adminOld).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "error": "Invalid Admin UUHash"})
		return
	}
	//Edit Target Admin Info
	var adminNew sqlTable.Admin = adminOld
	if body.AdminInfo.Password != "" {
		adminNew.PasswordHash = hash.SHA256(body.AdminInfo.Password)
	}
	if body.AdminInfo.Permission != adminOld.Permission {
		adminNew.Permission = body.AdminInfo.Permission
	}
	//Update Target Admin Info
	err = db.Save(&adminNew).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "error": "Internal Server Error"})
		return
	}
	//Return Success
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success"})
}
