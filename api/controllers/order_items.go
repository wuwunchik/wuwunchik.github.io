package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
)

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, order_id, dish_id, quantity FROM order_items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var oi models.OrderItem
		err := rows.Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orderItems = append(orderItems, oi)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderItems)
}

func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order item ID", http.StatusBadRequest)
		return
	}

	var oi models.OrderItem
	err = database.DB.QueryRow("SELECT id, order_id, dish_id, quantity FROM order_items WHERE id = ?", id).Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oi)
}

func AddDishToOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var orderItem models.OrderItem
	err = json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли уже такой пункт заказа
	var existingQuantity int
	err = database.DB.QueryRow("SELECT quantity FROM order_items WHERE order_id = ? AND dish_id = ?", orderID, orderItem.DishID).Scan(&existingQuantity)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err == nil {
		// Если пункт заказа уже существует, увеличиваем количество
		_, err = database.DB.Exec("UPDATE order_items SET quantity = quantity + ? WHERE order_id = ? AND dish_id = ?",
			orderItem.Quantity, orderID, orderItem.DishID)
	} else {
		// Если пункта заказа нет, создаем новый
		_, err = database.DB.Exec("INSERT INTO order_items (order_id, dish_id, quantity) VALUES (?, ?, ?)",
			orderID, orderItem.DishID, orderItem.Quantity)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderItem)
}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order item ID", http.StatusBadRequest)
		return
	}

	var oi models.OrderItem
	err = json.NewDecoder(r.Body).Decode(&oi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE order_items SET quantity = ? WHERE id = ?", oi.Quantity, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(oi)
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order item ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM order_items WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order item deleted successfully"})
}
