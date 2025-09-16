-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu;
DROP TABLE IF EXISTS dish_ingredients;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS units;

-- Создание таблицы единиц измерения
CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    abbreviation TEXT NOT NULL UNIQUE
);

-- Создание таблицы продуктов
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit_id INTEGER NOT NULL,
    FOREIGN KEY (unit_id) REFERENCES units(id)
);

-- Создание таблицы блюд
CREATE TABLE IF NOT EXISTS dishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL
);

-- Создание таблицы ингредиентов блюд
CREATE TABLE IF NOT EXISTS dish_ingredients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    UNIQUE (dish_id, product_id),
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Создание таблицы меню
CREATE TABLE IF NOT EXISTS menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    available BOOLEAN NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE
);

-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

-- Заполнение таблицы units начальными данными
INSERT INTO units (name, abbreviation) VALUES
('граммы', 'г'),
('килограммы', 'кг'),
('литры', 'л'),
('миллилитры', 'мл'),
('штуки', 'шт');

-- Заполнение таблицы products начальными данными
INSERT INTO products (name, quantity, unit_id) VALUES
('Мука', 100000, 2),   -- 100 кг
('Сахар', 50000, 2),   -- 50 кг
('Яйца', 200, 5),      -- 200 штук
('Молоко', 30000, 3),  -- 30 литров
('Масло', 20000, 2);   -- 20 кг

-- Заполнение таблицы dishes начальными данными
INSERT INTO dishes (name, description, price) VALUES
('Блинчики', 'С мукой и яйцами', 150.50),
('Суп', 'Овощной суп', 100.00),
('Котлета', 'Мясная котлета с гарниром', 200.00),
('Компот', 'Фруктовый компот', 50.00),
('Салат', 'Овощной салат', 80.00);

-- Заполнение таблицы dish_ingredients начальными данными
INSERT INTO dish_ingredients (dish_id, product_id, quantity) VALUES
(1, 1, 500),   -- Блинчики: 500 грамм муки
(1, 3, 10),    -- Блинчики: 10 яиц
(1, 4, 1000),  -- Блинчики: 1 литр молока (в миллилитрах)
(2, 2, 100),   -- Суп: 100 грамм сахара
(2, 4, 1000),  -- Суп: 1 литр молока (в миллилитрах)
(3, 1, 100),   -- Котлета: 100 грамм муки
(3, 5, 500);   -- Котлета: 500 грамм масла

-- Заполнение таблицы menu начальными данными
INSERT INTO menu (dish_id, available) VALUES
(1, 1),
(2, 1),
(3, 1),
(4, 1),
(5, 1);
