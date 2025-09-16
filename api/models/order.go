package models

import (
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	TableID   int       `json:"table_id"`
	OrderTime time.Time `json:"order_time"`
	Status    string    `json:"status"`
}

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"order_id"`
	DishID   int `json:"dish_id"`
	Quantity int `json:"quantity"`
}
