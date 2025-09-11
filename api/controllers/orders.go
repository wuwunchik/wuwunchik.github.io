package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.OrderTime = time.Now()

	// Проверяем, достаточно ли продуктов для заказа
	err = CheckProductAvailability(database.DB, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаём заказ
	_, err = database.DB.Exec("INSERT INTO orders (dish_id, quantity, order_time) VALUES (?, ?, ?)",
		order.DishID, order.Quantity, order.OrderTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Списание продуктов
	err = DeductProductsForOrder(database.DB, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func CheckProductAvailability(db *sql.DB, order models.Order) error {
	rows, err := db.Query(`
        SELECT p.id, p.quantity, di.quantity
        FROM products p
        JOIN dish_ingredients di ON p.id = di.product_id
        WHERE di.dish_id = ?
    `, order.DishID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, availableQuantity, requiredQuantity int
		if err := rows.Scan(&productID, &availableQuantity, &requiredQuantity); err != nil {
			return err
		}

		totalRequired := requiredQuantity * order.Quantity
		if availableQuantity < totalRequired {
			return fmt.Errorf("недостаточно продукта с ID %d", productID)
		}
	}

	return nil
}
