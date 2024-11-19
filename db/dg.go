package db

import (
	"database/sql"
	"log"
)

type Database struct {
	conn *sql.DB
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
	return &Database{conn: db}
}

func (db *Database) InitializeSchema()  {
	createExpensesTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT
		chat_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		category TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.conn.Exec(createExpensesTable)
	if err != nil {
		log.Fatalf("Create expenses table error: %v", err)
	}
	log.Println("Initialize schema successfully")
}

func (db *Database) AddExpenses(chatID int64, amount int, category string) error {
	query := `INSERT INTO expenses (chat_id, amount, category) (?, ?, ?)`
	_, err := db.conn.Exec(query, chatID, amount, category)
	if err != nil {
		return err
	}
	return nil
}
