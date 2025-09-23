-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu;
DROP TABLE IF EXISTS dish_ingredients;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS units;
DROP TABLE IF EXISTS tables;

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