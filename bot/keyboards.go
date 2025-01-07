package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMainKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить"),
			tgbotapi.NewKeyboardButton("Список"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Настройки"),
		),
	)
	return keyboard
}

func SettingsKb() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("День"),
			tgbotapi.NewKeyboardButton("Неделя"),
			tgbotapi.NewKeyboardButton("Месяц"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отмена"),
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
