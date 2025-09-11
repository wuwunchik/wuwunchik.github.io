-- Удаляем таблицы, если они существуют
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu;
DROP TABLE IF EXISTS dish_ingredients;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS products;

-- Создаём таблицу продуктов
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit TEXT NOT NULL
);

-- Создаём таблицу блюд
CREATE TABLE IF NOT EXISTS dishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL
);

-- Создаём таблицу ингредиентов блюд
CREATE TABLE IF NOT EXISTS dish_ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Создаём таблицу меню
CREATE TABLE IF NOT EXISTS menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    available BOOLEAN NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

-- Создаём таблицу заказов
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

-- Заполняем таблицу продуктов начальными данными
INSERT INTO products (name, quantity, unit) VALUES
('Мука', 100, 'кг'),
('Сахар', 50, 'кг'),
('Яйца', 200, 'шт'),
('Молоко', 30, 'л'),
('Масло', 20, 'кг');

-- Заполняем таблицу блюд начальными данными
INSERT INTO dishes (name, description, price) VALUES
('Блинчики', 'С мукой и яйцами', 150.50),
('Суп', 'Овощной суп', 100.00),
('Котлета', 'Мясная котлета с гарниром', 200.00),
('Компот', 'Фруктовый компот', 50.00),
('Салат', 'Овощной салат', 80.00);

-- Заполняем таблицу ингредиентов блюд начальными данными
INSERT INTO dish_ingredients (dish_id, product_id, quantity) VALUES
(1, 1, 500),   -- Блинчики: 500 грамм муки
(1, 3, 10),    -- Блинчики: 10 яиц
(1, 4, 1),     -- Блинчики: 1 литр молока
(2, 2, 100),   -- Суп: 100 грамм сахара
(2, 4, 1),     -- Суп: 1 литр молока
(3, 1, 100),   -- Котлета: 100 грамм муки
(3, 5, 500);   -- Котлета: 500 грамм масла

-- Заполняем таблицу меню начальными данными
INSERT INTO menu (dish_id, available) VALUES
(1, 1),
(2, 1),
(3, 1),
(4, 1),
(5, 1);
