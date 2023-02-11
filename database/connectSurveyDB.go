package database

import (
	"fmt"
	"log"

	"github.com/drkgrntt/duffy-json-api/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	surveyDB *gorm.DB
)

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

func GetSurveyDatabase() *gorm.DB {
	return surveyDB
}
