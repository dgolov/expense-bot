package db

import (
	"fmt"
)

type Expense struct {
	Amount			int
	Category  		string
	CreatedAt       string
}

type Budget struct {
	Amount			int
	Spent   		int
}


func (exp *Expense) GetText() string {
	return fmt.Sprintf("%d руб. (%s)", exp.Amount, exp.CreatedAt)
}

func (exp *Expense) GetTextWithCategory() string {
	return fmt.Sprintf("%d руб. на %s (%s)", exp.Amount, exp.Category, exp.CreatedAt)
}
