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

func GetDishes(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, description, price FROM dishes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dishes []models.Dish
	for rows.Next() {
		var d models.Dish
		err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dishes)
}

func GetDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid dish ID", http.StatusBadRequest)
		return
	}

	var d models.Dish
	err = database.DB.QueryRow("SELECT id, name, description, price FROM dishes WHERE id = ?", id).Scan(&d.ID, &d.Name, &d.Description, &d.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Dish not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO dishes (name, description, price) VALUES (?, ?, ?)", d.Name, d.Description, d.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	d.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func UpdateDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid dish ID", http.StatusBadRequest)
		return
	}

	var d models.Dish
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE dishes SET name = ?, description = ?, price = ? WHERE id = ?", d.Name, d.Description, d.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func DeleteDish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid dish ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM dishes WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
