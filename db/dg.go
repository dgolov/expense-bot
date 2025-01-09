package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	Conn *sql.DB
}

func NewDatabase(dataSourceName string) *Database {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Connect to database error: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Ping database error: %v", err)
	}

	log.Println("Connect to database successfully")
	return &Database{Conn: db}
}

func (db *Database) InitializeSchema() {
	go db.createExpenses()
	go db.createBudgets()
}

func (db *Database) createExpenses() {
	log.Println("Start createExpenses function")
	createExpensesTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		category TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Conn.Exec(createExpensesTable)
	if err != nil {
		log.Fatalf("Create expenses table error: %v", err)
	}
	log.Println("Initialize expenses schema successfully")
}

func (db *Database) createBudgets() {
	log.Println("Start createBudgets function")
	createBudgetsTable := `
	CREATE TABLE IF NOT EXISTS budget (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		budget_amount INTEGER NOT NULL,
		spent_amount INTEGER DEFAULT 0
	);
	`

	_, err := db.Conn.Exec(createBudgetsTable)
	if err != nil {
		log.Fatalf("Create budgets table error: %v", err)
	}
	log.Println("Initialize budgets schema successfully")
}

func (db *Database) AddExpenses(chatID int64, amount int, category string) error {
	query := `INSERT INTO expenses (chat_id, amount, category) VALUES  (?, ?, ?)`
	_, err := db.Conn.Exec(query, chatID, amount, category)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ListExpenses(chatID int64, period string) ([]*Expense, error) {
	var query string

	switch period {
	case "day":
		query = `
			SELECT category, SUM(amount), MAX(created_at) 
			FROM expenses 
			WHERE chat_id = ? 
				AND strftime('%Y-%m-%d', created_at) = strftime('%Y-%m-%d', 'now')
			GROUP BY category
			ORDER BY MAX(created_at) DESC
		`
	case "week":
		query = `
			SELECT category, SUM(amount), MAX(created_at) 
			FROM expenses 
			WHERE chat_id = ? 
				AND created_at >= datetime('now', '-7 days')
			GROUP BY category
			ORDER BY MAX(created_at) DESC
		`
	case "month":
		query = `
			SELECT category, SUM(amount), MAX(created_at) 
			FROM expenses 
			WHERE chat_id = ? 
				AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')
			GROUP BY category
			ORDER BY MAX(created_at) DESC
		`
	default:
		return nil, fmt.Errorf("Invalid period: %s", period)
	}

	rows, err := db.Conn.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return execListExpensesQuery(rows)
}

func (db *Database) ListExpensesByCategory(chatID int64, category string, period string) ([]*Expense, error) {
	var query string

	switch period {
	case "day":
		query = `
			SELECT category, amount, created_at
			FROM expenses 
			WHERE chat_id = ? 
				AND strftime('%Y-%m-%d', created_at) = strftime('%Y-%m-%d', 'now')
				AND category = ?
			ORDER BY created_at DESC
		`
	case "week":
		query = `
			SELECT category, amount, created_at
			FROM expenses 
			WHERE chat_id = ? 
				AND created_at >= datetime('now', '-7 days')
				AND category = ?
			ORDER BY created_at DESC
		`
	case "month":
		query = `
			SELECT category, amount, created_at
			FROM expenses 
			WHERE chat_id = ? 
				AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')
				AND category = ?
			ORDER BY created_at DESC
		`
	default:
		return nil, fmt.Errorf("Invalid period: %s", period)
	}

	rows, err := db.Conn.Query(query, chatID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return execListExpensesQuery(rows)
}

func execListExpensesQuery(rows *sql.Rows) ([]*Expense, error) {
	var expenses []*Expense

	for rows.Next() {
		var category string
		var amount int
		var createdAt string

		err := rows.Scan(&category, &amount, &createdAt)
		if err != nil {
			return nil, err
		}

		expense := &Expense{
			Amount: amount,
			Category: category,
			CreatedAt: createdAt,
		}

		expenses = append(expenses, expense)
	}
	return expenses, nil
}
