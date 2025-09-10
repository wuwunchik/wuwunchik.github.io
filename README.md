# Cafeteria Management System

**Простое API для управления столовой: продукты, блюда, меню и заказы.**

---

## 📋 Описание проекта

Этот проект реализует бэкенд для управления столовой с использованием **SQLite** и **Go**. API поддерживает CRUD-операции для:

- **Продуктов** (складской учёт)
- **Блюд** (меню столовой)
- **Меню** (доступность блюд)
- **Заказов** (история заказов)

Проект включает:

- Готовые SQL-скрипты для создания таблиц.
- CRUD-операции на Go.
- Интеграцию с **Swagger UI** для документации и тестирования API.

---

## 🛠 Структура базы данных

### Таблицы

| Таблица    | Описание                       |
| ---------- | ------------------------------ |
| `products` | Продукты на складе             |
| `dishes`   | Блюда, доступные в столовой    |
| `menu`     | Текущее меню (связь с блюдами) |
| `orders`   | История заказов                |

### Схема базы данных

```sql
-- Продукты
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit TEXT NOT NULL
);

-- Блюда
CREATE TABLE dishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL
);

-- Меню (связь блюд и их доступности)
CREATE TABLE menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    available BOOLEAN NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

-- Журнал заказов
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dish_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);
```
