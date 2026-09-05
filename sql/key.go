package sqlOperator

import (
	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

func AddApplicationKey(db *gorm.DB, keyRow sqlTable.ApplicationKey) error {
	return db.Create(&keyRow).Error
}

func GetApplicationKey(db *gorm.DB, keyRow string) (sqlTable.ApplicationKey, error) {
	var key sqlTable.ApplicationKey
	if err := db.Where("key = ?", keyRow).First(&key).Error; err != nil {
		return key, err
	}
	return key, nil
}

func CheckApplicationKey(db *gorm.DB, keyRow string) bool {
	var key sqlTable.ApplicationKey
	if err := db.Where("key = ?", keyRow).First(&key).Error; err != nil {
		return false
	}
	return true
}

func DeleteApplicationKey(db *gorm.DB, key string) error {
	return db.Where("key = ?", key).Delete(&sqlTable.ApplicationKey{}).Error
}
