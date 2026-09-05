package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetCharacterByCharacterUUHash(db *gorm.DB, characterUUHash string) (sqlTable.Character, error) {
	var character sqlTable.Character
	err := db.Where("uu_hash = ?", characterUUHash).First(&character).Error
	return character, err
}

func GetCharacterByUserUUHash(db *gorm.DB, userUUHash string) ([]sqlTable.Character, error) {
	var characters []sqlTable.Character
	err := db.Where("user_uu_hash = ?", userUUHash).Find(&characters).Error
	return characters, err
}
