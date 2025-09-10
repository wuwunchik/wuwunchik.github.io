package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/products.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы Products (если не существует)
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected and migrated!")
}
