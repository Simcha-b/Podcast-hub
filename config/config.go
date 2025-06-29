package config

import (
	"os"

	"github.com/Simcha-b/Podcast-Hub/utils"
	"github.com/joho/godotenv"
)

var Logger = utils.NewLogger("info")

type config struct {
	PORT          string
	DATA_DIR      string
	TIME_INTERVAL string
}

func LoadConfig() *config {
	err := godotenv.Load()
	if err != nil {
		Logger.Error("Error loading .env file")
		os.Exit(1)
	}
	config := config{
		PORT:          os.Getenv("PORT"),
		DATA_DIR:      os.Getenv("DATA_DIR"),
		TIME_INTERVAL: os.Getenv("TIME_INTERVAL"),
	}
	return &config
}
