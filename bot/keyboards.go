package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить"),
			tgbotapi.NewKeyboardButton("Список"),
		),
	)
	return keyboard
}

func GetCancelKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
		),
	)
	return keyboard
}
