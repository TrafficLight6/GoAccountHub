package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func AddCharacter(db *gorm.DB, character sqlTable.Character) error {
	return db.Create(&character).Error
}
