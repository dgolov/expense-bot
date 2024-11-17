package bot

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API				  *tgbotapi.BotAPI
	AwaitingExpenses  map[int64]bool
	Timers            map[int64]*time.Timer
}

func (b *Bot) SetAwaitingExpense(chatID int64)  {
	b.AwaitingExpenses[chatID] = true

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
	}

	b.Timers[chatID] = time.AfterFunc(5 * time.Minute, func() {
		b.ResetAwaitingExpense(chatID)
		msg := tgbotapi.NewMessage(chatID, "Время ожидания истекло. Попробуйте снова отправить команду /add.")
		b.API.Send(msg)
	})
}

func (b *Bot) ResetAwaitingExpense(chatID int64) {
	log.Printf("ResetAwaitingExpense for %d", chatID)

	delete(b.AwaitingExpenses, chatID)

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
		delete(b.Timers, chatID)
	}
}

func NewBot(botToken string, debugMode bool) *Bot  {
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = debugMode
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	return &Bot{
		API:			  botAPI,
		AwaitingExpenses: make(map[int64]bool),
		Timers:           make(map[int64]*time.Timer),
	}
}
