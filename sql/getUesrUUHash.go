package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetUserUUHash(db *gorm.DB, userID uint) (string, error) {
	var user sqlTable.User
	err := db.Where("id = ?", userID).First(&user).Error
	return user.UUHash, err
}
