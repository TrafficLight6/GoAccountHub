// If Admin is Root Admin,its id is 0
package sqlTable

import (
	"time"

	"gorm.io/gorm"
)

type AdminLoginToken struct {
	gorm.Model
	ID       uint      `gorm:"primaryKey"`
	UUHash   string    `gorm:"column:uu_hash"`
	Token    string    `gorm:"column:token"`
	DeadTime time.Time `gorm:"column:dead_time"`
}
