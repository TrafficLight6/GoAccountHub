package sqlOperator

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TrafficLight6/GoAccountHub/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("Error in Connecting DB:", err)
		return nil
	}
	InitTable(db)
	return db
}
