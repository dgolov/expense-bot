package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DebugMode bool
}

func LoadConfig() *Config {
	log.Println("Loading config")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Telegram token does not found in env file")
	}

	debugStr := os.Getenv("DEBUG")
	debugMode := false
	if debugStr != "" {
		debugMode, err = strconv.ParseBool(debugStr)
		if err != nil {
			log.Printf("Parsing DEBUG error: %v. Default value: false", err)
		}
		log.Printf("Debug: %t", debugMode)
	}

	return &Config{
		BotToken: botToken,
		DebugMode: debugMode,
	}
}
