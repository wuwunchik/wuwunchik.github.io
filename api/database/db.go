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

	// Создание таблицы столиков
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS tables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			number INTEGER NOT NULL UNIQUE,
			capacity INTEGER NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы блюд
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS dishes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы ингредиентов блюд
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS dish_ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dish_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			UNIQUE (dish_id, product_id),
			FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы меню
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS menu (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dish_id INTEGER NOT NULL,
			available BOOLEAN NOT NULL,
			FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы заказов
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			table_id INTEGER NOT NULL,
			order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT NOT NULL DEFAULT 'created',
			FOREIGN KEY (table_id) REFERENCES tables(id)
		);
	`)
	if err != nil {
		return err
	}

	// Создание таблицы пунктов заказа
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			dish_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			UNIQUE (order_id, dish_id),
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			FOREIGN KEY (dish_id) REFERENCES dishes(id)
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

	// Заполнение таблицы tables начальными данными
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO tables (number, capacity) VALUES
		(1, 4),
		(2, 4),
		(3, 6),
		(4, 2),
		(5, 8);
	`)
	if err != nil {
		return err
	}

	// Заполнение таблицы dishes начальными данными
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO dishes (name, description, price) VALUES
		('Блинчики', 'С мукой и яйцами', 150.50),
		('Суп', 'Овощной суп', 100.00),
		('Котлета', 'Мясная котлета с гарниром', 200.00),
		('Компот', 'Фруктовый компот', 50.00),
		('Салат', 'Овощной салат', 80.00);
	`)
	if err != nil {
		return err
	}

	// Заполнение таблицы dish_ingredients начальными данными
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO dish_ingredients (dish_id, product_id, quantity) VALUES
		(1, 1, 500),   -- Блинчики: 500 грамм муки
		(1, 3, 10),    -- Блинчики: 10 яиц
		(1, 4, 1000),  -- Блинчики: 1 литр молока (в миллилитрах)
		(2, 2, 100),   -- Суп: 100 грамм сахара
		(2, 4, 1000),  -- Суп: 1 литр молока (в миллилитрах)
		(3, 1, 100),   -- Котлета: 100 грамм муки
		(3, 5, 500);   -- Котлета: 500 грамм масла
	`)
	if err != nil {
		return err
	}

	// Заполнение таблицы menu начальными данными
	_, err = DB.Exec(`
		INSERT OR IGNORE INTO menu (dish_id, available) VALUES
		(1, 1),
		(2, 1),
		(3, 1),
		(4, 1),
		(5, 1);
	`)
	if err != nil {
		return err
	}

	return nil
}
