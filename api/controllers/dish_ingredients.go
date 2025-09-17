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

func GetDishIngredients(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, dish_id, product_id, quantity FROM dish_ingredients")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var dishIngredients []models.DishIngredient
	for rows.Next() {
		var di models.DishIngredient
		err := rows.Scan(&di.ID, &di.DishID, &di.ProductID, &di.Quantity)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dishIngredients = append(dishIngredients, di)
	}

	utils.RespondWithJSON(w, http.StatusOK, dishIngredients)
}

func GetDishIngredient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ingredient ID")
		return
	}

	var di models.DishIngredient
	err = database.DB.QueryRow("SELECT id, dish_id, product_id, quantity FROM dish_ingredients WHERE id = ?", id).Scan(&di.ID, &di.DishID, &di.ProductID, &di.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Dish ingredient not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, di)
}

func CreateDishIngredient(w http.ResponseWriter, r *http.Request) {
	var di models.DishIngredient
	err := json.NewDecoder(r.Body).Decode(&di)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("INSERT INTO dish_ingredients (dish_id, product_id, quantity) VALUES (?, ?, ?)", di.DishID, di.ProductID, di.Quantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(di)
}

func UpdateDishIngredient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ingredient ID")
		return
	}

	var di models.DishIngredient
	err = json.NewDecoder(r.Body).Decode(&di)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE dish_ingredients SET dish_id = ?, product_id = ?, quantity = ? WHERE id = ?", di.DishID, di.ProductID, di.Quantity, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, di)
}

func DeleteDishIngredient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ingredient ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM dish_ingredients WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Dish ingredient deleted successfully"})
}
