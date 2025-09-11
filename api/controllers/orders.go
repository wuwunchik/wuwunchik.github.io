package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
)

// Функция для проверки доступности продуктов
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

// Функция для списания продуктов при заказе
func DeductProductsForOrder(db *sql.DB, order models.Order) error {
	rows, err := db.Query(`
		SELECT product_id, quantity
		FROM dish_ingredients
		WHERE dish_id = ?
	`, order.DishID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, requiredQuantity int
		if err := rows.Scan(&productID, &requiredQuantity); err != nil {
			return err
		}

		totalRequired := requiredQuantity * order.Quantity

		_, err = db.Exec(`
			UPDATE products
			SET quantity = quantity - ?
			WHERE id = ? AND quantity >= ?
		`, totalRequired, productID, totalRequired)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, dish_id, quantity, order_time FROM orders")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.DishID, &o.Quantity, &o.OrderTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var o models.Order
	err = database.DB.QueryRow("SELECT id, dish_id, quantity, order_time FROM orders WHERE id = ?", id).Scan(&o.ID, &o.DishID, &o.Quantity, &o.OrderTime)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.OrderTime = time.Now()

	err = CheckProductAvailability(database.DB, o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO orders (dish_id, quantity, order_time) VALUES (?, ?, ?)", o.DishID, o.Quantity, o.OrderTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = DeductProductsForOrder(database.DB, o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var o models.Order
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE orders SET dish_id = ?, quantity = ? WHERE id = ?", o.DishID, o.Quantity, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
}
