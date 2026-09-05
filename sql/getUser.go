package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetUserByUsernameAndPassword(db *gorm.DB, username string, passwordHash string) (sqlTable.User, error) {
	var user sqlTable.User
	err := db.Where("username = ? AND password_hash = ?", username, passwordHash).First(&user).Error
	return user, err
}
