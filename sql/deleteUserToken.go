package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func DeleteUserToken(db *gorm.DB, token string) error {
	if err := db.Where("token = ?", token).Delete(&sqlTable.UserLoginToken{}).Error; err != nil {
		return err
	}
	return nil
}
