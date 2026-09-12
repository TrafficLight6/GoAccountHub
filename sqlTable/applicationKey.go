package sqlTable

import "gorm.io/gorm"

type ApplicationKey struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`
	//UNIQUE
	KeyUser string `gorm:"column:key_user"`
	Key     string `gorm:"column:key"`
}
