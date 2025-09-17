package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
	"wuwunchik.github.io/api/utils"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, table_id, order_time, status FROM orders")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.TableID, &o.OrderTime, &o.Status)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, o)
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var orderResponse models.OrderResponse

	// Получаем информацию о заказе
	err = database.DB.QueryRow(`
		SELECT id, table_id, order_time, status
		FROM orders
		WHERE id = ?
	`, id).Scan(&orderResponse.ID, &orderResponse.Table.ID, &orderResponse.OrderTime, &orderResponse.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Получаем пункты заказа с информацией о блюдах
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.order_id, oi.dish_id, oi.quantity, d.id, d.name, d.description, d.price
		FROM order_items oi
		JOIN dishes d ON oi.dish_id = d.id
		WHERE oi.order_id = ?
	`, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var orderItemWithDish models.OrderItemWithDish
		var dish models.Dish

		err := rows.Scan(
			&orderItemWithDish.ID,
			&orderItemWithDish.OrderID,
			&orderItemWithDish.Dish.ID,
			&orderItemWithDish.Quantity,
			&dish.ID,
			&dish.Name,
			&dish.Description,
			&dish.Price)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		orderItemWithDish.Dish = dish
		orderResponse.Table.Dishes = append(orderResponse.Table.Dishes, orderItemWithDish)
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]models.OrderResponse{"order": orderResponse})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	o.OrderTime = time.Now()
	o.Status = "created"

	_, err = database.DB.Exec("INSERT INTO orders (table_id, order_time, status) VALUES (?, ?, ?)",
		o.TableID, o.OrderTime, o.Status)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, o)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var o models.Order
	err = json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE orders SET table_id = ?, status = ? WHERE id = ?", o.TableID, o.Status, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, o)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	// Начинаем транзакцию
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// Получаем все пункты заказа перед удалением заказа
	rows, err := tx.Query("SELECT id, dish_id, quantity FROM order_items WHERE order_id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	// Возвращаем продукты на склад для всех пунктов заказа
	for rows.Next() {
		var oi models.OrderItem
		err := rows.Scan(&oi.ID, &oi.DishID, &oi.Quantity)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		oi.OrderID = id
		err = ReturnProductsForOrder(tx, oi)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Удаляем заказ
	_, err = tx.Exec("DELETE FROM orders WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}
