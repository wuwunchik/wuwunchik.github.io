# Cafeteria Management System

**API для управления столовой: продукты, блюда, меню и заказы.**

---

## 📋 Описание проекта

API поддерживает CRUD-операции для:

- **Продуктов** (складской учёт)
- **Блюд** (меню столовой)
- **Меню** (доступность блюд)
- **Заказов** (история заказов)

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

## 🚀 Быстрый старт

1. Клонируйте репозиторий

```bash
git clone https://github.com/wuwunchik/wuwunchik.github.io.git
cd wuwunchik.github.io

```

2. Установите зависимости

```bash

go mod tidy
```

3. Запустите сервер

```bash

go run main.go
```

Сервер будет доступен по адресу: http://localhost:8080.

## 📚 Документация API

Документация API доступна через Swagger UI: [Cafeteria Management API](https://wuwunchik.github.io/api/v_0_1/#/)

## 📂 Структура проекта

```text

📦wuwunchik.github.io
├── 📂docs
│ └── 📂api
│ └── 📂v_0_1
│ ├── swagger.yaml
│ └── index.html
├── 📜main.go
├── 📜go.mod
└── 📜README.md

```
