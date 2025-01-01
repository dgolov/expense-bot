package main

import (
	"expense-bot/bot"
	"expense-bot/config"
	"expense-bot/db"
	"log"
)

func main()  {
	log.Println("Bot is started")

	cfg := config.LoadConfig()

	database := db.NewDatabase("expenses.db")
	defer database.Conn.Close()
	database.InitializeSchema()

	expenseBot := bot.NewBot(cfg.BotToken, cfg.DebugMode, cfg.TimeoutMinutes, database)
	expenseBot.HandleUpdates()
}
