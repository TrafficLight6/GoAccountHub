package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

// Using Hard Delete To Delete Admin Token From Database
func DelAdminToken(db *gorm.DB, token string) error {
	//Hard Delete: Tokens Are Ephemeral Data, No Need to Keep Soft Deleted Rows
	return db.Unscoped().Where("token = ?", token).Delete(&sqlTable.AdminLoginToken{}).Error
}
