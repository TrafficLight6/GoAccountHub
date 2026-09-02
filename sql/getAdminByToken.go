package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetAdminByToken(db *gorm.DB, token string, appConfig config.Config) (sqlTable.Admin, error) {
	var adminLoginToken sqlTable.AdminLoginToken
	if err := db.Where("token = ?", token).First(&adminLoginToken).Error; err != nil {
		return sqlTable.Admin{}, err
	}
	//Check is Root Admin
	if adminLoginToken.UUHash == appConfig.RootAdminUUHash {
		return sqlTable.Admin{Username: "root"}, nil
	}
	//Check is Normal Admin
	var admin sqlTable.Admin
	if err := db.Where("uu_hash = ?", adminLoginToken.UUHash).First(&admin).Error; err != nil {
		return sqlTable.Admin{}, err
	}
	return admin, nil
}
