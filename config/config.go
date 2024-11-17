package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken		  string
	DebugMode 		  bool
	TimeoutMinutes    int
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
		} else {
			log.Printf("Debug: %t", debugMode)
		}
	}

	timeoutStr := os.Getenv("TIMEOUT_MINUTES")
	timeoutMinutes := 5
	if timeoutStr != "" {
		timeoutMinutes, err = strconv.Atoi(timeoutStr)
		if err != nil {
			log.Printf("Parsing TIMEOUT_MINUTES error: %v. Default value: 5", err)
		} else {
			log.Printf("TIMEOUT_MINUTES: %v", timeoutMinutes)
		}
	}

	return &Config{
		BotToken: botToken,
		DebugMode: debugMode,
		TimeoutMinutes: timeoutMinutes,
	}
}
