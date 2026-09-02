package sqlTable

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	UUHash       string `gorm:"column:uu_hash"`

	Permission Permission `gorm:"column:permission;type:jsonb;serializer:json"`
}

type Permission struct {
}
