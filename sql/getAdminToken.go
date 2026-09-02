package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetAdminToken(db *gorm.DB, token string) (sqlTable.AdminLoginToken, error) {
	var admin sqlTable.AdminLoginToken
	err := db.Where("token = ?", token).First(&admin).Error
	return admin, err
}
