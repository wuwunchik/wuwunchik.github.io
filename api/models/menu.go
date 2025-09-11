package models

type MenuItem struct {
	ID        int  `json:"id"`
	DishID    int  `json:"dish_id"`
	Available bool `json:"available"`
}
