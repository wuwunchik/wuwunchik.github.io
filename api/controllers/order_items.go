package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
	"wuwunchik.github.io/api/utils"
)

// Функция для проверки доступности продуктов
func CheckProductAvailability(db *sql.DB, orderItem models.OrderItem) error {
	rows, err := db.Query(`
		SELECT p.id, p.quantity, di.quantity
		FROM products p
		JOIN dish_ingredients di ON p.id = di.product_id
		WHERE di.dish_id = ?
	`, orderItem.DishID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, availableQuantity, requiredQuantity int
		if err := rows.Scan(&productID, &availableQuantity, &requiredQuantity); err != nil {
			return err
		}

		totalRequired := requiredQuantity * orderItem.Quantity
		if availableQuantity < totalRequired {
			return fmt.Errorf("недостаточно продукта с ID %d", productID)
		}
	}

	return nil
}

// Функция для списания продуктов при заказе
func DeductProductsForOrder(db *sql.DB, orderItem models.OrderItem) error {
	rows, err := db.Query(`
		SELECT product_id, quantity
		FROM dish_ingredients
		WHERE dish_id = ?
	`, orderItem.DishID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, requiredQuantity int
		if err := rows.Scan(&productID, &requiredQuantity); err != nil {
			return err
		}

		totalRequired := requiredQuantity * orderItem.Quantity
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

// Функция для возврата продуктов при уменьшении количества в заказе
func ReturnProductsForOrder(db *sql.DB, orderItem models.OrderItem) error {
	rows, err := db.Query(`
		SELECT product_id, quantity
		FROM dish_ingredients
		WHERE dish_id = ?
	`, orderItem.DishID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, requiredQuantity int
		if err := rows.Scan(&productID, &requiredQuantity); err != nil {
			return err
		}

		totalReturned := requiredQuantity * (-orderItem.Quantity) // Отрицательное значение для возврата
		_, err = db.Exec(`
			UPDATE products
			SET quantity = quantity + ?
			WHERE id = ?
		`, totalReturned, productID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, order_id, dish_id, quantity FROM order_items")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var oi models.OrderItem
		err := rows.Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orderItems = append(orderItems, oi)
	}

	utils.RespondWithJSON(w, http.StatusOK, orderItems)
}

func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	var oi models.OrderItem
	err = database.DB.QueryRow("SELECT id, order_id, dish_id, quantity FROM order_items WHERE id = ?", id).Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order item not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, oi)
}

func AddDishToOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var orderItem models.OrderItem
	err = json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderItem.OrderID = orderID

	// Проверяем доступность продуктов
	err = CheckProductAvailability(database.DB, orderItem)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем, существует ли уже такой пункт заказа
	var existingQuantity int
	err = database.DB.QueryRow("SELECT quantity FROM order_items WHERE order_id = ? AND dish_id = ?", orderItem.OrderID, orderItem.DishID).Scan(&existingQuantity)
	if err != nil && err != sql.ErrNoRows {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err == nil {
		// Если пункт заказа уже существует, увеличиваем количество
		_, err = database.DB.Exec("UPDATE order_items SET quantity = quantity + ? WHERE order_id = ? AND dish_id = ?",
			orderItem.Quantity, orderItem.OrderID, orderItem.DishID)
	} else {
		// Если пункта заказа нет, создаем новый
		_, err = database.DB.Exec("INSERT INTO order_items (order_id, dish_id, quantity) VALUES (?, ?, ?)",
			orderItem.OrderID, orderItem.DishID, orderItem.Quantity)
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Списание продуктов
	err = DeductProductsForOrder(database.DB, orderItem)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, orderItem)
}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	var oi models.OrderItem
	err = json.NewDecoder(r.Body).Decode(&oi)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем текущее количество для расчета разницы
	var currentQuantity int
	err = database.DB.QueryRow("SELECT quantity FROM order_items WHERE id = ?", id).Scan(&currentQuantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Обновляем количество в пункте заказа
	_, err = database.DB.Exec("UPDATE order_items SET quantity = ? WHERE id = ?", oi.Quantity, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Корректируем количество продуктов
	diff := oi.Quantity - currentQuantity
	if diff != 0 {
		var orderItem models.OrderItem
		orderItem.DishID = oi.DishID
		orderItem.Quantity = diff

		if diff > 0 {
			// Проверяем доступность продуктов
			err = CheckProductAvailability(database.DB, orderItem)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			// Списание продуктов
			err = DeductProductsForOrder(database.DB, orderItem)
		} else {
			// Возврат продуктов
			err = ReturnProductsForOrder(database.DB, orderItem)
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, oi)
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid order item ID")
		return
	}

	// Получаем информацию о пункте заказа перед удалением
	var oi models.OrderItem
	err = database.DB.QueryRow("SELECT id, order_id, dish_id, quantity FROM order_items WHERE id = ?", id).Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Order item not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Удаляем пункт заказа
	_, err = database.DB.Exec("DELETE FROM order_items WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Возвращаем продукты на склад
	err = ReturnProductsForOrder(database.DB, oi)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Order item deleted successfully"})
}
