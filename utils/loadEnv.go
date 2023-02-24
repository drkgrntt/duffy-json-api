package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Version     string `mapstructure:"VERSION"`
	Environment string `mapstructure:"ENVIRONMENT"`

	ServerPort string `mapstructure:"PORT"`
	IpInfoKey  string `mapstructure:"IP_INFO_KEY"`

	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`

	SurveyDBHost     string `mapstructure:"MYSQL_HOST"`
	SurveyDBUser     string `mapstructure:"MYSQL_USER"`
	SurveyDBPassword string `mapstructure:"MYSQL_PASSWORD"`
	SurveyDBName     string `mapstructure:"MYSQL_DB"`
	SurveyDBPort     string `mapstructure:"MYSQL_PORT"`

	HpfDBHost     string `mapstructure:MONGO_HOST`
	HpfDBUser     string `mapstructure:MONGO_USER`
	HpfDBPassword string `mapstructure:MONGO_PASSWORD`
	HpfDBName     string `mapstructure:MONGO_DB`
}

var ConfigInstance Config

func LoadConfig(path string) (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	local := godotenv.Overload("local.env")
	if local != nil {
		log.Println("No local config found")
	}

	config = Config{
		Version:     os.Getenv("VERSION"),
		Environment: os.Getenv("ENVIRONMENT"),
		ServerPort:  os.Getenv("PORT"),
		IpInfoKey:   os.Getenv("IP_INFO_KEY"),

		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBName:     os.Getenv("POSTGRES_DB"),
		DBPort:     os.Getenv("POSTGRES_PORT"),

		SurveyDBHost:     os.Getenv("MYSQL_HOST"),
		SurveyDBUser:     os.Getenv("MYSQL_USER"),
		SurveyDBPassword: os.Getenv("MYSQL_PASSWORD"),
		SurveyDBName:     os.Getenv("MYSQL_DB"),
		SurveyDBPort:     os.Getenv("MYSQL_PORT"),

		HpfDBHost:     os.Getenv("MONGO_HOST"),
		HpfDBUser:     os.Getenv("MONGO_USER"),
		HpfDBPassword: os.Getenv("MONGO_PASSWORD"),
		HpfDBName:     os.Getenv("MONGO_DB"),
	}

	ConfigInstance = config
	return
}

func GetConfig() Config {
	return ConfigInstance
}
