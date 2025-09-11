package models

import (
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	DishID    int       `json:"dish_id"`
	Quantity  int       `json:"quantity"`
	OrderTime time.Time `json:"order_time"`
}
