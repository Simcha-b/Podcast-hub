package config

import "os"

type config struct {
	PORT     string
	DATA_DIR string
}

func LoadConfig() *config {
	config := config{
		PORT:     os.Getenv("PORT"),
		DATA_DIR: os.Getenv("DATA_DIR"),
	}
	return &config
}
