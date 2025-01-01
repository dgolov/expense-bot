package bot

import (
	"expense-bot/db"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdates(storage *db.Database)  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				start(b, chatID)
			case "add":
				add(b, chatID)
			case "list":
				list(b, chatID, storage)
			case "cancel":
				cancel(b, chatID)
			default:
				if checkMessage(update.Message.Text, b, chatID, storage) == 1 {
					continue
				}
				if b.AwaitingExpenses[chatID] {
					if strings.Contains(update.Message.Text, " ") {
						parts := strings.SplitN(update.Message.Text, " ", 2)
						amount, err := strconv.Atoi(parts[0])
						if err != nil {
							msg := tgbotapi.NewMessage(chatID,
								"Ошибка: сумма должна быть числом.\nЕсли передумали, отправьте /cancel.")
							b.API.Send(msg)
							continue
						}
						category := parts[1]
						err = storage.AddExpenses(chatID, amount, category)
						if err != nil {
							log.Printf("Add expenses error: %v", err)
							b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка добавления расходов."))
							return
						} else {
							b.API.Send(tgbotapi.NewMessage(chatID, "Расход добавлен!"))
						}

						b.ResetAwaitingExpense(chatID)
					} else {
						msg := tgbotapi.NewMessage(chatID,
							"Ошибка: введите расход в формате <сумма> <категория>.")
						b.API.Send(msg)
					}
				} else {
					b.API.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда."))
				}
			}
		}
	}
}

func checkMessage(text string, b *Bot, chatID int64, storage *db.Database) int8 {
	switch text {
	case "Добавить":
		add(b, chatID)
		return 1
	case "Список":
		list(b, chatID, storage)
		return 1
	case "Отмена":
		cancel(b, chatID)
		return 1
	}
	return 0
}

func start(b *Bot, chatID int64) {
	keyboard := getMainKb()
	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Я помогу вам вести учет расходов.")
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func add(b *Bot, chatID int64)  {
	keyboard := getCancelKb()
	b.SetAwaitingExpense(chatID)
	msg := tgbotapi.NewMessage(chatID,
		"Введите расход в формате: <сумма> <категория>.\nЕсли передумали, отправьте /cancel.")
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func list(b *Bot, chatID int64, storage *db.Database) {
	keyboard := getMainKb()
	expenses, err := storage.ListExpenses(chatID)
	if err != nil {
		log.Printf("Get expenses error: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}
	if len(expenses) == 0 {
		msg := tgbotapi.NewMessage(chatID, "У вас пока нет расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Ваши расходы:\n" + strings.Join(expenses, "\n"))
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

func cancel(b *Bot, chatID int64) {
	keyboard := getMainKb()
	if b.AwaitingExpenses[chatID] {
		b.ResetAwaitingExpense(chatID)
		b.AwaitingExpenses[chatID] = false
		msg := tgbotapi.NewMessage(chatID, "Добавление расхода отменено.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Вы не находитесь в процессе добавления расхода.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

func getMainKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить"),
			tgbotapi.NewKeyboardButton("Список"),
		),
	)
	return keyboard
}

func getCancelKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
	return keyboard
}
