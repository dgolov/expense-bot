package main

import (
	"expense-bot/bot"
	"expense-bot/config"
	"expense-bot/models"
	"log"
)

func main()  {
	log.Println("Bot is started")

	cfg := config.LoadConfig()
	expenseBot := bot.NewBot(cfg.BotToken, cfg.DebugMode, cfg.TimeoutMinutes)
	storage := models.NewStorage()

	expenseBot.HandleUpdates(storage)
}
