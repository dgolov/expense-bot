package models

import "fmt"

type Storage struct {
	Expenses map[int64][]string
}

func NewStorage() *Storage  {
	return &Storage{Expenses: make(map[int64][]string)}
}

func (s *Storage) AddExpense(chatID int64, amount int, category string) {
	expense := fmt.Sprintf("%d %s", amount, category)
	s.Expenses[chatID] = append(s.Expenses[chatID], expense)
}

func (s *Storage) ListExpenses(chatID int64) []string {
	return s.Expenses[chatID]
}
