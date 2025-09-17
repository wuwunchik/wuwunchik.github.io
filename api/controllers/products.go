package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/models"
	"wuwunchik.github.io/api/utils"
)

// GetProducts возвращает список всех продуктов
func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT p.id, p.name, p.quantity, p.unit_id, u.id, u.name, u.abbreviation
		FROM products p
		JOIN units u ON p.unit_id = u.id
	`)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var unit models.Unit
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.UnitID, &unit.ID, &unit.Name, &unit.Abbreviation)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.Unit = unit
		products = append(products, p)
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

// GetProduct возвращает информацию о продукте
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p models.Product
	var unit models.Unit
	err = database.DB.QueryRow(`
		SELECT p.id, p.name, p.quantity, p.unit_id, u.id, u.name, u.abbreviation
		FROM products p
		JOIN units u ON p.unit_id = u.id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.Name, &p.Quantity, &p.UnitID, &unit.ID, &unit.Name, &unit.Abbreviation)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	p.Unit = unit

	utils.RespondWithJSON(w, http.StatusOK, p)
}

// CreateProduct создает новый продукт
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.DB.Exec("INSERT INTO products (name, quantity, unit_id) VALUES (?, ?, ?)", p.Name, p.Quantity, p.UnitID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)

	// Получаем информацию о единице измерения для возврата полной информации о продукте
	var unit models.Unit
	err = database.DB.QueryRow("SELECT id, name, abbreviation FROM units WHERE id = ?", p.UnitID).Scan(&unit.ID, &unit.Name, &unit.Abbreviation)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	p.Unit = unit

	utils.RespondWithJSON(w, http.StatusCreated, p)
}

// UpdateProduct обновляет информацию о продукте
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p models.Product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE products SET name = ?, quantity = ?, unit_id = ? WHERE id = ?", p.Name, p.Quantity, p.UnitID, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	p.ID = id

	// Получаем информацию о единице измерения для возврата полной информации о продукте
	var unit models.Unit
	err = database.DB.QueryRow("SELECT id, name, abbreviation FROM units WHERE id = ?", p.UnitID).Scan(&unit.ID, &unit.Name, &unit.Abbreviation)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	p.Unit = unit

	utils.RespondWithJSON(w, http.StatusOK, p)
}

// DeleteProduct удаляет продукт
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
