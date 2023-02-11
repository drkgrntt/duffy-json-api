package database

import (
	"fmt"
	"log"

	"github.com/drkgrntt/duffy-json-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func ConnectDB(config *utils.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	fmt.Println("? Connecting using the following DSN: " + dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Failed to connect to the PG Database")
	}
	fmt.Println("? Connected Successfully to the PG Database")
}

func GetDatabase() *gorm.DB {
	return DB
}
