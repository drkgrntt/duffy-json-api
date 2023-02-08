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

	DBHost     string `mapstructure:"POSTGRES_HOST"`
	DBUser     string `mapstructure:"POSTGRES_USER"`
	DBPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName     string `mapstructure:"POSTGRES_DB"`
	DBPort     string `mapstructure:"POSTGRES_PORT"`
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
		DBHost:      os.Getenv("POSTGRES_HOST"),
		DBUser:      os.Getenv("POSTGRES_USER"),
		DBPassword:  os.Getenv("POSTGRES_PASSWORD"),
		DBName:      os.Getenv("POSTGRES_DB"),
		DBPort:      os.Getenv("POSTGRES_PORT"),
	}

	ConfigInstance = config
	return
}

func GetConfig() Config {
	return ConfigInstance
}
