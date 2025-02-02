package bot

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleStart(b *Bot, chatID int64) {
	log.Println("Handler start")
	keyboard := GetMainKb()
	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Я помогу вам вести учет расходов.")
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleAdd(b *Bot, chatID int64)  {
	log.Println("Handler add")
	keyboard := GetCancelKb()
	b.SetAwaitingExpense(chatID)
	msgText := "Введите расход в формате: <сумма> <категория>.\nЕсли передумали, отправьте /cancel."
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleSettings(b *Bot, chatID int64)  {
	log.Println("Handler settings")
	keyboard := SettingsKb()
	b.SetAwaitingSettings(chatID)

	msgText := "\nВы вошли в режим настроек.\n"
	msgText += "\nУкажите период за который нужно учитывать расходы.\nУстановите или измените ваш бюджет.\n"
	msgText += "\nЕсли передумали и хотите выйти, отправьте /cancel."
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleSetPeriod(b *Bot, chatID int64, period string)  {
	var msgText string
	log.Println("Handler set period")
	keyboard := GetMainKb()

	if b.AwaitingSettings[chatID] {
		msgText = "Выбран период учета расходов - " + period + "."
		b.SetPeriod(period)
	} else {
		msgText = "Вы не находитесь в режиме настроек.\nДля перехода, отправьте /settings."
	}
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)

	b.ResetAwaitingSettings(chatID)
}

func handleSave(b *Bot, text string, chatID int64) {
	log.Println("Handler save")
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

func getBudget(b *Bot, chatID int64) {
	log.Println("Handler get budget")
	keyboard := GetMainKb()
	var msgTxt string

	budget, err := b.Storage.GetBudgetByChatId(chatID)
	if err != nil {
		log.Printf("Get budget error: %v", err)
		msgTxt = "Ошибка при получении бюджета."
	} else {
		log.Printf("Get budget: %v", budget)
		if budget == nil {
			msgTxt = "Бюджет не установлен."
		} else {
			msgTxt = "Бюджет: " + strconv.Itoa(budget.Amount) + "\nПотрачено: "  + strconv.Itoa(budget.Spent)
			if budget.Spent > budget.Amount {
				msgTxt += "\nВнимание! Ваши расходы превышают бюджет!"
			}
		}
	}

	msg := tgbotapi.NewMessage(chatID, msgTxt)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func setBudget(b *Bot, chatID int64, text string) {
	log.Printf("Handler set budget. Text - %s", text)

	var msgTxt string
	var amount int
	var err error

	if text == "" {
		keyboard := GetCancelKb()
		b.SetAwaitingBudget(chatID)
		msgText := "Укажите сумму бюджета."
		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	keyboard := GetMainKb()

	if b.AwaitingBudget[chatID] {
		b.ResetAwaitingBudget(chatID)
		b.AwaitingExpenses[chatID] = false
		amount, err = strconv.Atoi(text)
	} else {
		parts := strings.SplitN(text, " ", 3)
		amount, err = strconv.Atoi(parts[2])
		log.Printf("Handler set budget. Amount - %d", amount)
	}

	if err != nil {
		msgText := "Ошибка: сумма бюджета должна быть числом."
		msgTxt += "\nЕсли передумали, отправьте /cancel."
		msg := tgbotapi.NewMessage(chatID, msgText)
		b.API.Send(msg)
		return
	}

	err = b.Storage.SetBudgetForChatId(chatID, amount)
	if err != nil {
		log.Printf("Set budget error: %v", err)
		msgTxt = "Ошибка установки бюджета."
	} else {
		log.Println("Set budget successfully")
		msgTxt = "Бюджет успешно установлен."
	}

	msg := tgbotapi.NewMessage(chatID, msgTxt)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func handleList(b *Bot, chatID int64) {
	log.Println("Handler list")
	keyboard := GetMainKb()

	expenses, err := b.Storage.ListExpenses(chatID, b.Period)
	if err != nil {
		log.Printf("Get expenses error: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	period := TranslatePeriod(b.Period)
	if len(expenses) == 0 {
		msg := tgbotapi.NewMessage(chatID, "За " + period + " у вас пока нет расходов.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		var expensesTxtList []string
		for _, itemExpense  := range expenses {
			expensesTxtList = append(expensesTxtList, itemExpense.GetTextWithCategory())
		}

		msgTxt := "Ваши расходы за " + period + ":\n" + strings.Join(expensesTxtList, "\n")
		msg := tgbotapi.NewMessage(chatID, msgTxt)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

func handleListByCategory(b *Bot, chatID int64, text string) {
	log.Println("Handler list by category")
	parts := strings.SplitN(text, " ", 2)
	category := parts[1]
	keyboard := GetMainKb()
	expenses, err := b.Storage.ListExpensesByCategory(chatID, category, b.Period)
	if err != nil {
		log.Printf("Get expenses error: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при получении расходов по категории " + category + ".")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
		return
	}

	period := TranslatePeriod(b.Period)
	if len(expenses) == 0 {
		msgTxt := "За " + period + " у вас пока нет расходов по категории" + category + "."
		msg := tgbotapi.NewMessage(chatID, msgTxt)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		var expensesTxtList []string
		for _, itemExpense  := range expenses {
			expensesTxtList = append(expensesTxtList, itemExpense.GetText())
		}
		msgTxt := "Ваши расходы за " + period + " на " + category + ":\n"
		msgTxt = msgTxt + strings.Join(expensesTxtList, "\n")
		msg := tgbotapi.NewMessage(chatID, msgTxt)
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

func handleCancel(b *Bot, chatID int64) {
	log.Println("Handler cancel")
	keyboard := GetMainKb()
	if b.AwaitingExpenses[chatID] {
		b.ResetAwaitingExpense(chatID)
		b.AwaitingExpenses[chatID] = false
		msg := tgbotapi.NewMessage(chatID, "Добавление расхода отменено.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else if b.AwaitingSettings[chatID]  {
		b.ResetAwaitingSettings(chatID)
		b.AwaitingSettings[chatID] = false
		msg := tgbotapi.NewMessage(chatID, "Выход из режима настроек.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Вы не находитесь в процессе добавления расхода.")
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}
