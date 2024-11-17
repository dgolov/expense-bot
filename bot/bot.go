package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API *tgbotapi.BotAPI
	AwaitingExpenses  map[int64]bool
}

func NewBot(botToken string) *Bot  {
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	return &Bot{
		API: botAPI,
		AwaitingExpenses: make(map[int64]bool),
	}
}
