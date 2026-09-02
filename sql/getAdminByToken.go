package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetAdminByToken(db *gorm.DB, token string) (sqlTable.Admin, error) {
	var adminToken sqlTable.AdminLoginToken
	if err := db.Where("token = ?", token).First(&adminToken).Error; err != nil {
		return sqlTable.Admin{}, err
	}
	//Check is Root Admin
	if adminToken.AdminID == 0 {
		return sqlTable.Admin{Username: "root"}, nil
	}
	//Check is Normal Admin
	var admin sqlTable.Admin
	if err := db.Where("id = ?", adminToken.AdminID).First(&admin).Error; err != nil {
		return sqlTable.Admin{}, err
	}
	return admin, nil
}
