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

	// Создание таблицы ролей
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS roles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            description TEXT
        );
    `)
	if err != nil {
		return err
	}

	// Создание таблицы пользователей
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		return err
	}

	// Создание таблицы пересечений roles и users
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS user_roles (
			user_id INTEGER NOT NULL,
			role_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
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

	// Добавляем роли
	_, err = DB.Exec(`
	INSERT OR IGNORE INTO roles (id, name, description) VALUES
	(1, 'admin', 'Администратор системы'),
	(2, 'user', 'Обычный пользователь'),
	(3, 'manager', 'Менеджер склада');
`)
	if err != nil {
		return err
	}

	// Добавляем админа с ролью admin
	adminPasswordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMy.MrYv5eLmQ5Q9QqZ4Z7J3J5J1J0J4J7J" // хеш для "password"
	_, err = DB.Exec(`
	INSERT OR IGNORE INTO users (id, username, password_hash) VALUES
	(1, 'admin', ?);
`, adminPasswordHash)
	if err != nil {
		return err
	}

	// Назначаем админу роль admin
	_, err = DB.Exec(`
	INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES
	(1, 1);  -- user_id=1 (admin), role_id=1 (admin)
`)
	if err != nil {
		return err
	}

	return nil
}
