package sqlTable

import "gorm.io/gorm"

type ApplicationKey struct {
	gorm.Model
	ID                   int64  `gorm:"primaryKey"`
	ApplicationUsingSite string `gorm:"column:application_using_site"`
	Key                  string `gorm:"column:key"`
}
