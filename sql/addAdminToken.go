package sqlOperator

import (
	"time"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func AddAdminToken(db *gorm.DB, adminUUHash string, token string, deadTime time.Time) error {
	loginToken := sqlTable.AdminLoginToken{
		UUHash:   adminUUHash,
		Token:    token,
		DeadTime: deadTime,
	}
	return db.Create(&loginToken).Error
}
