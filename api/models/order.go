package models

import (
	"time"
)

// Order структура для базовой информации о заказе
type Order struct {
	ID        int       `json:"id"`
	TableID   int       `json:"table_id"`
	OrderTime time.Time `json:"order_time"`
	Status    string    `json:"status"`
}

// OrderItem структура для базовой информации о пункте заказа
type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"order_id"`
	DishID   int `json:"dish_id"`
	Quantity int `json:"quantity"`
}

// OrderResponse структура для детализированного ответа о заказе
type OrderResponse struct {
	ID        int             `json:"id"`
	Table     TableWithDishes `json:"table_id"`
	OrderTime string          `json:"order_time"`
	Status    string          `json:"status"`
}

// TableWithDishes структура для информации о столике с блюдами
type TableWithDishes struct {
	ID     int                 `json:"id"`
	Dishes []OrderItemWithDish `json:"dishes"`
}

// OrderItemWithDish структура для информации о пункте заказа с блюдом
type OrderItemWithDish struct {
	ID       int  `json:"id"`
	OrderID  int  `json:"order_id"`
	Dish     Dish `json:"dish"`
	Quantity int  `json:"quantity"`
}
