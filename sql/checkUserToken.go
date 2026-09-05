package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func CheckUserToken(db *gorm.DB, token string) bool {
	var tokenRecord sqlTable.UserLoginToken
	if err := db.Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return false
	}
	return true
}
