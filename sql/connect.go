package sqlOperator

import (
	"fmt"

	"github.com/TrafficLight6/GoAccountHub/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(appConfig config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DatabaseHost,
		appConfig.DatabasePort,
		appConfig.DatabaseUser,
		appConfig.DatabasePassword,
		appConfig.DatabaseName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error in Connecting DB:", err)
		return nil
	}
	InitTable(db)
	return db
}
