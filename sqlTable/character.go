package sqlTable

import "gorm.io/gorm"

type Character struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	CharacterName string `gorm:"column:character_name"`
	PasswordHash  string `gorm:"column:password_hash"`
	UserUUHash    string `gorm:"column:user_uu_hash"`
	UUHash        string `gorm:"column:uu_hash"`
	//Can be Empty,or Json, XML, Yaml, and Else
	MetaData string `gorm:"column:meta_data"`
}
