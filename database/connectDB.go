package database

import (
	"fmt"
	"log"

	"github.com/drkgrntt/duffy-json-api/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB       *gorm.DB
	surveyDB *gorm.DB
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

func ConnectSurveyDB(config *utils.Config) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.SurveyDBUser, config.SurveyDBPassword, config.SurveyDBHost, config.SurveyDBPort, config.SurveyDBName)

	fmt.Println("? Connecting using the following DSN: " + dsn)
	surveyDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Failed to connect to the Survey Database")
	}
	fmt.Println("? Connected Successfully to the Survey Database")
}

func GetDatabase() *gorm.DB {
	return DB
}

func GetSurveyDatabase() *gorm.DB {
	return surveyDB
}
