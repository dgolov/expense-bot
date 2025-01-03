package bot

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleStart(b *Bot, chatID int64) {
	keyboard := GetMainKb()
	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Я помогу вам вести учет расходов.")
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleAdd(b *Bot, chatID int64)  {
	keyboard := GetCancelKb()
	b.SetAwaitingExpense(chatID)
	msgText := "Введите расход в формате: <сумма> <категория>.\nЕсли передумали, отправьте /cancel."
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleSave(b *Bot, text string, chatID int64) {
	parts := strings.SplitN(text, " ", 2)
	amount, err := strconv.Atoi(parts[0])
	if err != nil {
		msgText := "Ошибка: сумма должна быть числом.\nЕсли передумали, отправьте /cancel."
		msg := tgbotapi.NewMessage(chatID, msgText)
		b.API.Send(msg)
		return
	}

	keyboard := GetMainKb()

	category := parts[1]
	err = b.Storage.AddExpenses(chatID, amount, category)
	if err != nil {
		log.Printf("Add expenses error: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка добавления расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	} else {
		msg := tgbotapi.NewMessage(chatID, "Расход добавлен!")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}

	b.ResetAwaitingExpense(chatID)
}

func handleList(b *Bot, chatID int64) {
	keyboard := GetMainKb()
	expenses, err := b.Storage.ListExpenses(chatID)
	if err != nil {
		log.Printf("Get expenses error: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}
	if len(expenses) == 0 {
		msg := tgbotapi.NewMessage(chatID, "За текущий месяцу у вас пока нет расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Ваши расходы за текущий месяц:\n" + strings.Join(expenses, "\n"))
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

func handleCancel(b *Bot, chatID int64) {
	keyboard := GetMainKb()
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
