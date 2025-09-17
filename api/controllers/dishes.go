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

func GetDishes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, price FROM dishes")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var dishes []models.Dish
	for rows.Next() {
		var d models.Dish
		err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.Price)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dishes = append(dishes, d)
	}

	utils.RespondWithJSON(w, http.StatusOK, dishes)
}

func GetDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ID")
		return
	}

	var d models.Dish
	err = database.DB.QueryRow("SELECT id, name, description, price FROM dishes WHERE id = ?", id).Scan(&d.ID, &d.Name, &d.Description, &d.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Dish not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, d)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := database.DB.Exec("INSERT INTO dishes (name, description, price) VALUES (?, ?, ?)", d.Name, d.Description, d.Price)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	d.ID = int(id)

	utils.RespondWithJSON(w, http.StatusCreated, d)
}

func UpdateDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ID")
		return
	}

	var d models.Dish
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE dishes SET name = ?, description = ?, price = ? WHERE id = ?", d.Name, d.Description, d.Price, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	d.ID = id
	utils.RespondWithJSON(w, http.StatusOK, d)
}

func DeleteDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid dish ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM dishes WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Dish deleted successfully"})
}
