package models

type DishIngredient struct {
	ID        int `json:"id"`
	DishID    int `json:"dish_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
