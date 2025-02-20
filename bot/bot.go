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
	AwaitingSettings  map[int64]bool
	AwaitingBudget    map[int64]bool
	Timers            map[int64]*time.Timer
	TimeoutMinutes    int
	Storage 	 	  *db.Database
	Period 			  string
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
		AwaitingSettings: make(map[int64]bool),
		AwaitingBudget:   make(map[int64]bool),
		Timers:           make(map[int64]*time.Timer),
		TimeoutMinutes:   timeoutMinutes,
		Storage:   		  storage,
		Period: 		  "month",
	}
}

func (b *Bot) SetAwaitingExpense(chatID int64)  {
	log.Printf("SetAwaitingExpense for %d", chatID)
	b.AwaitingExpenses[chatID] = true

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
	}

	timeoutDuration := time.Duration(b.TimeoutMinutes) * time.Minute
	b.Timers[chatID] = time.AfterFunc(timeoutDuration, func() {
		b.ResetAwaitingExpense(chatID)
		keyboard := GetMainKb()
		msg := tgbotapi.NewMessage(chatID, "Время ожидания истекло. Попробуйте снова отправить команду /add.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	})
}

func (b *Bot) SetAwaitingSettings(chatID int64)  {
	log.Printf("SetAwaitingSettings for %d", chatID)
	b.AwaitingSettings[chatID] = true

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
	}

	timeoutDuration := time.Duration(b.TimeoutMinutes) * time.Minute
	b.Timers[chatID] = time.AfterFunc(timeoutDuration, func() {
		b.ResetAwaitingSettings(chatID)
		keyboard := GetMainKb()
		msg := tgbotapi.NewMessage(chatID, "Время ожидания истекло. Попробуйте снова отправить команду /settings.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	})
}

func (b *Bot) SetAwaitingBudget(chatID int64)  {
	log.Printf("SetAwaitingBudget for %d", chatID)
	b.AwaitingBudget[chatID] = true

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
	}

	timeoutDuration := time.Duration(b.TimeoutMinutes) * time.Minute
	b.Timers[chatID] = time.AfterFunc(timeoutDuration, func() {
		b.ResetAwaitingExpense(chatID)
		msg := tgbotapi.NewMessage(chatID, "Время ожидания истекло. Попробуйте снова отправить команду /set_budget.")
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

func (b *Bot) ResetAwaitingSettings(chatID int64) {
	log.Printf("ResetAwaitingSettings for %d", chatID)

	delete(b.AwaitingSettings, chatID)

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
		delete(b.Timers, chatID)
	}
}

func (b *Bot) ResetAwaitingBudget(chatID int64) {
	log.Printf("ResetAwaitingBudget for %d", chatID)

	delete(b.AwaitingBudget, chatID)

	if timer, exists := b.Timers[chatID]; exists {
		timer.Stop()
		delete(b.Timers, chatID)
	}
}

func (b *Bot) SetPeriod(period string)  {
	log.Printf("Set period - %s", period)
	b.Period = period
}

func (b *Bot) HandleUpdates()  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("User: [%s] - %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				handleStart(b, chatID)
			case "add":
				handleAdd(b, chatID)
			case "list":
				handleList(b, chatID)
			case "settings":
				handleSettings(b, chatID)
			case "day":
				handleSetPeriod(b, chatID, "day")
			case "week":
				handleSetPeriod(b, chatID, "week")
			case "month":
				handleSetPeriod(b, chatID, "month")
			case "cancel":
				handleCancel(b, chatID)
			case "budget":
				getBudget(b, chatID)
			case "set_budget":
				setBudget(b, chatID, "")
			default:
				if b.checkMessage(update.Message.Text, chatID) == 1 {
					continue
				}
				if msg := strings.ToLower(update.Message.Text);
					strings.Contains(msg, "установи бюджет") ||
						strings.Contains(msg, "установить бюджет") {
					setBudget(b, chatID, update.Message.Text)
					continue
				}
				if b.AwaitingBudget[chatID] {
					setBudget(b, chatID, update.Message.Text)
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
	text = strings.ToLower(text)

	switch text {
	case "добавить":
		handleAdd(b, chatID)
		return 1
	case "список":
		handleList(b, chatID)
		return 1
	case "отмена":
		handleCancel(b, chatID)
		return 1
	case "настройки":
		handleSettings(b, chatID)
		return 1
	case "день":
		handleSetPeriod(b, chatID, "day")
		return 1
	case "неделя":
		handleSetPeriod(b, chatID, "week")
		return 1
	case "месяц":
		handleSetPeriod(b, chatID, "month")
		return 1
	case "бюджет":
		getBudget(b, chatID)
		return 1
	case "установить бюджет":
		setBudget(b, chatID, "")
		return 1
	}

	if strings.Contains(text, "список") {
		handleListByCategory(b, chatID, text)
		return 1
	}

	return 0
}
