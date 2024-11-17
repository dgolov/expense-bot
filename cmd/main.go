package main

import (
	"expense-bot/bot"
	"expense-bot/config"
	"expense-bot/models"
)

func main()  {
	cfg := config.LoadConfig()
	expenseBot := bot.NewBot(cfg.BotToken)
	storage := models.NewStorage()

	expenseBot.HandleUpdates(storage)
}
