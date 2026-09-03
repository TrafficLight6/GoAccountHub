package sqlTable

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	UUHash       string `gorm:"column:uu_hash"` //CONSTANT!!
	//Can be Empty,or Json, XML, Yaml, and Else
	MetaData string `gorm:"column:meta_data"`
}
