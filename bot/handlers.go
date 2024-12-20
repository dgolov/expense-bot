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
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Я помогу вам вести учет расходов.")
				b.API.Send(msg)

			case "add":
				b.SetAwaitingExpense(chatID)
				msg := tgbotapi.NewMessage(chatID,
					"Введите расход в формате: <сумма> <категория>.\nЕсли передумали, отправьте /cancel.")
				b.API.Send(msg)

			case "list":
				expenses, err := storage.ListExpenses(chatID)
				if err != nil {
					log.Printf("Get expenses error: %v", err)
					b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении расходов."))
					return
				}
				if len(expenses) == 0 {
					b.API.Send(tgbotapi.NewMessage(chatID, "У вас пока нет расходов."))
				} else {
					msg := "Ваши расходы:\n" + strings.Join(expenses, "\n")
					b.API.Send(tgbotapi.NewMessage(chatID, msg))
				}

			case "cancel":
				if b.AwaitingExpenses[chatID] {
					b.ResetAwaitingExpense(chatID)
					b.AwaitingExpenses[chatID] = false
					msg := tgbotapi.NewMessage(chatID, "Добавление расхода отменено.")
					b.API.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Вы не находитесь в процессе добавления расхода.")
					b.API.Send(msg)
				}

			default:
				if b.AwaitingExpenses[chatID] {
					if strings.Contains(update.Message.Text, " ") {
						parts := strings.SplitN(update.Message.Text, "", 2)
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
