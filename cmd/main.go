package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
)


var expenses = make(map[int64][]string)


func main()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID,
					"Добро пожаловать! Я помогу вам вести учет расходов.\n" +
					"Используйте команду /add для добавления расхода.")
				bot.Send(msg)

			case "add":
				msg := tgbotapi.NewMessage(chatID,
					"Введите расход в формате: <сумма> <категория>\nПример: 500 еда")
				bot.Send(msg)

			case "":
				if strings.Contains(update.Message.Text, " ") {
					parts := strings.SplitN(update.Message.Text ,"", 2)
					ammout, err := strconv.Atoi(parts[0])
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Ошибка: сумма должна быть числом."))
						continue
					}
					category := parts[1]
					expense := fmt.Sprintf("%d %s", ammout, category)
					expenses[chatID] = append(expenses[chatID], expense)
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка: введите расход в правильном формате."))
				}

			case "list":
				if userExpenses, exits := expenses[chatID]; exits {
					msg := "Ваши расходы:\n" + strings.Join(userExpenses, "\n")
					bot.Send(tgbotapi.NewMessage(chatID, msg))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, "У вас пока нет расходов."))
				}

			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда. Попробуйте /start, /add или /list."))
			}
		}
	}
}
