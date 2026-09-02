package sqlOperator

import (
	"time"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func InsertAdminToken(db *gorm.DB, adminID uint, token string, deadTime time.Time) error {
	loginToken := sqlTable.AdminLoginToken{
		AdminID:  adminID,
		Token:    token,
		DeadTime: deadTime,
	}
	return db.Create(&loginToken).Error
}
