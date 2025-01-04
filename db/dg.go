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

func (db *Database) InitializeSchema()  {
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
	log.Println("Initialize schema successfully")
}

func (db *Database) AddExpenses(chatID int64, amount int, category string) error {
	query := `INSERT INTO expenses (chat_id, amount, category) VALUES  (?, ?, ?)`
	_, err := db.Conn.Exec(query, chatID, amount, category)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ListExpenses(chatID int64) ([]string, error) {
	query := `
		SELECT category, SUM(amount), MAX(created_at) 
		FROM expenses 
		WHERE chat_id = ? 
			AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')
		GROUP BY category
		ORDER BY MAX(created_at) DESC
	`
	rows, err := db.Conn.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []string
	for rows.Next() {
		var category string
		var totalAmount int
		var lastUpdated string

		err := rows.Scan(&category, &totalAmount, &lastUpdated)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, fmt.Sprintf("%d руб. на %s (%s)", totalAmount, category, lastUpdated))
	}
	return expenses, nil
}

func (db *Database) ListExpensesByCategory(chatID int64, category string) ([]string, error) {
	query := `
		SELECT category, amount, created_at
		FROM expenses 
		WHERE chat_id = ? 
			AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')
			AND category = ?
		ORDER BY created_at DESC
	`
	rows, err := db.Conn.Query(query, chatID, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []string
	for rows.Next() {
		var category string
		var totalAmount int
		var lastUpdated string

		err := rows.Scan(&category, &totalAmount, &lastUpdated)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, fmt.Sprintf("%d руб. (%s)", totalAmount, lastUpdated))
	}
	return expenses, nil
}
