package bot

import (
	"expense-bot/models"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdates(storage *models.Storage)  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Я помогу вам вести учет расходов.")
				b.API.Send(msg)

			case "add":
				msg := tgbotapi.NewMessage(chatID, "Введите расход в формате: <сумма> <категория>")
				b.API.Send(msg)

			case "list":
				expenses := storage.ListExpenses(chatID)
				if len(expenses) == 0 {
					b.API.Send(tgbotapi.NewMessage(chatID, "У вас пока нет расходов."))
				} else {
					msg := "Ваши расходы:\n" + strings.Join(expenses, "\n")
					b.API.Send(tgbotapi.NewMessage(chatID, msg))
				}

			default:
				if strings.Contains(update.Message.Text, " ") {
					parts := strings.SplitN(update.Message.Text ,"", 2)
					amount, err := strconv.Atoi(parts[0])
					if err != nil {
						b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка: сумма должна быть числом."))
						continue
					}
					category := parts[1]
					storage.AddExpense(chatID, amount, category)
					b.API.Send(tgbotapi.NewMessage(chatID, "Расход добавлен!"))
				} else {
					b.API.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда."))
				}
			}
		}
	}
}
