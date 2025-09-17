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

func GetMenu(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, dish_id, available FROM menu")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var menuItems []models.MenuItem
	for rows.Next() {
		var mi models.MenuItem
		err := rows.Scan(&mi.ID, &mi.DishID, &mi.Available)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		menuItems = append(menuItems, mi)
	}

	utils.RespondWithJSON(w, http.StatusOK, menuItems)
}

func GetMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item ID")
		return
	}

	var mi models.MenuItem
	err = database.DB.QueryRow("SELECT id, dish_id, available FROM menu WHERE id = ?", id).Scan(&mi.ID, &mi.DishID, &mi.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Menu item not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, mi)
}

func CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var mi models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&mi)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("INSERT INTO menu (dish_id, available) VALUES (?, ?)", mi.DishID, mi.Available)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, mi)
}

func UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item ID")
		return
	}

	var mi models.MenuItem
	err = json.NewDecoder(r.Body).Decode(&mi)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE menu SET dish_id = ?, available = ? WHERE id = ?", mi.DishID, mi.Available, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, mi)
}

func DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM menu WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}
