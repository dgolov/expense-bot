package bot

import (
	"expense-bot/db"
	"log"
	"time"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API				  *tgbotapi.BotAPI
	AwaitingExpenses  map[int64]bool
	Timers            map[int64]*time.Timer
	TimeoutMinutes    int
	Storage 	 	  *db.Database
}

func NewBot(botToken string, debugMode bool, timeoutMinutes int, storage *db.Database) *Bot  {
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
		TimeoutMinutes:   timeoutMinutes,
		Storage:   		  storage,
	}
}

func (b *Bot) SetAwaitingExpense(chatID int64)  {
	b.AwaitingExpenses[chatID] = true

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
	}

	timeoutDuration := time.Duration(b.TimeoutMinutes) * time.Minute
	b.Timers[chatID] = time.AfterFunc(timeoutDuration, func() {
		b.ResetAwaitingExpense(chatID)
		msg := tgbotapi.NewMessage(chatID, "Время ожидания истекло. Попробуйте снова отправить команду /add.")
		b.API.Send(msg)
	})
}

func (b *Bot) SaveExpensesToDB(text string, chatID int64) {
	if strings.Contains(text, " ") {
		handleSave(b, text, chatID)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Ошибка: введите расход в формате <сумма> <категория>.")
		b.API.Send(msg)
	}
}

func (b *Bot) ResetAwaitingExpense(chatID int64) {
	log.Printf("ResetAwaitingExpense for %d", chatID)

	delete(b.AwaitingExpenses, chatID)

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
		delete(b.Timers, chatID)
	}
}

func (b *Bot) HandleUpdates()  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				handleStart(b, chatID)
			case "add":
				handleAdd(b, chatID)
			case "list":
				handleList(b, chatID)
			case "cancel":
				handleCancel(b, chatID)
			default:
				if b.checkMessage(update.Message.Text, chatID) == 1 {
					continue
				}
				if b.AwaitingExpenses[chatID] {
					b.SaveExpensesToDB(update.Message.Text, chatID)
				} else {
					b.API.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда."))
				}
			}
		}
	}
}

func (b *Bot) checkMessage(text string, chatID int64) int8 {
	switch text {
	case "Добавить":
		handleAdd(b, chatID)
		return 1
	case "Список":
		handleList(b, chatID)
		return 1
	case "Отмена":
		handleCancel(b, chatID)
		return 1
	}
	return 0
}
