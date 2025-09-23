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
	DB, err = sql.Open("sqlite3", "./store.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблиц
	err = createTables()
	if err != nil {
		log.Fatal(err)
	}

	// Заполнение таблиц начальными данными
	err = seedTables()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully!")
}

func createTables() error {
	// Создание таблицы единиц измерения
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS units (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			abbreviation TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы продуктов
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit_id INTEGER NOT NULL,
			FOREIGN KEY (unit_id) REFERENCES units(id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func seedTables() error {
	// Заполнение таблицы units начальными данными
	_, err := DB.Exec(`
		INSERT OR IGNORE INTO units (name, abbreviation) VALUES
		('граммы', 'г'),
		('килограммы', 'кг'),
		('литры', 'л'),
		('миллилитры', 'мл'),
		('штуки', 'шт');
	`)
	if err != nil {
		return err
	}

	// Заполнение таблицы products начальными данными
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO products (name, quantity, unit_id) VALUES
		('Мука', 100000, 2),   -- 100 кг
		('Сахар', 50000, 2),   -- 50 кг
		('Яйца', 200, 5),      -- 200 штук
		('Молоко', 30000, 3),  -- 30 литров
		('Масло', 20000, 2);   -- 20 кг
	`)
	if err != nil {
		return err
	}

	return nil
}
