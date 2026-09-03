package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

// Check Admin Token Exist
func CheckAdminToken(db *gorm.DB, token string) bool {
	var tokenRecord sqlTable.AdminLoginToken
	if err := db.Where("token = ?", token).First(&tokenRecord).Error; err != nil {
		return false
	}
	return true
}
