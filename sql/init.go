package sqlOperator

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/sqlTable"
	"gorm.io/gorm"
)

var (
	//Tables to be created
	Tables = []any{
		&sqlTable.Admin{},
		&sqlTable.AdminLoginToken{},

		&sqlTable.User{},
		&sqlTable.UserLoginToken{},

		&sqlTable.Character{},

		&sqlTable.ApplicationKey{},
	}
)

func InitTable(db *gorm.DB) {
	for _, table := range Tables {
		if db.Migrator().HasTable(table) {
			continue
		}
		//Create Table
		db.AutoMigrate(table)
		if err := db.AutoMigrate(table); err != nil {
			fmt.Println("Cannot Create Table:", err)
		}
	}
}
