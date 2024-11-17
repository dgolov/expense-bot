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
	expenseBot := bot.NewBot(cfg.BotToken, cfg.DebugMode)
	storage := models.NewStorage()

	expenseBot.HandleUpdates(storage)
}
