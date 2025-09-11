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

func GetMenu(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, dish_id, available FROM menu")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var menuItems []models.MenuItem
	for rows.Next() {
		var mi models.MenuItem
		err := rows.Scan(&mi.ID, &mi.DishID, &mi.Available)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		menuItems = append(menuItems, mi)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menuItems)
}

func GetMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	var mi models.MenuItem
	err = database.DB.QueryRow("SELECT id, dish_id, available FROM menu WHERE id = ?", id).Scan(&mi.ID, &mi.DishID, &mi.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Menu item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mi)
}

func CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var mi models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&mi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO menu (dish_id, available) VALUES (?, ?)", mi.DishID, mi.Available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mi)
}

func UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	var mi models.MenuItem
	err = json.NewDecoder(r.Body).Decode(&mi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE menu SET dish_id = ?, available = ? WHERE id = ?", mi.DishID, mi.Available, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mi)
}

func DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM menu WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Menu item deleted successfully"})
}
