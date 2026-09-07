package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func GetCharacterByToken(db *gorm.DB, token string) (sqlTable.Character, error) {
	var userLoginToken sqlTable.UserLoginToken
	if err := db.Where("token = ?", token).First(&userLoginToken).Error; err != nil {
		return sqlTable.Character{}, err
	}
	var character sqlTable.Character
	if err := db.Where("uu_hash = ? AND user_uu_hash = ?", userLoginToken.CharacterUUHash, userLoginToken.UUHash).First(&character).Error; err != nil {
		return sqlTable.Character{}, err
	}
	return character, nil
}
