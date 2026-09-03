package sqlTable

import (
	"time"

	"gorm.io/gorm"
)

type UserLoginToken struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey"`
	UUHash          string    `gorm:"column:uu_hash"`
	CharacterUUHash string    `gorm:"column:character_uu_hash"`
	Token           string    `gorm:"column:token"`
	DeadTime        time.Time `gorm:"column:dead_time"`
}
