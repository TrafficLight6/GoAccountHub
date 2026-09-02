package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetAdminByUsernameAndPassword(db *gorm.DB, username string, passwordHash string) (sqlTable.Admin, error) {
	var admin sqlTable.Admin
	err := db.Where("username = ? AND password_hash = ?", username, passwordHash).First(&admin).Error
	return admin, err
}
