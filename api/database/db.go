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
	DB, err = sql.Open("sqlite3", "./cafeteria.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблиц
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS dishes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS menu (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dish_id INTEGER NOT NULL,
			available BOOLEAN NOT NULL,
			FOREIGN KEY (dish_id) REFERENCES dishes(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dish_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (dish_id) REFERENCES dishes(id)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully!")
}
