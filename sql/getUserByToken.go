package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetUserByToken(db *gorm.DB, token string) (sqlTable.User, error) {
	var userLoginToken sqlTable.UserLoginToken
	if err := db.Where("token = ?", token).First(&userLoginToken).Error; err != nil {
		return sqlTable.User{}, err
	}
	var user sqlTable.User
	if err := db.Where("uu_hash = ?", userLoginToken.UUHash).First(&user).Error; err != nil {
		return sqlTable.User{}, err
	}
	return user, nil
}
