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

func GetTables(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, number, capacity FROM tables")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(&t.ID, &t.Number, &t.Capacity)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tables = append(tables, t)
	}

	utils.RespondWithJSON(w, http.StatusOK, tables)
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	var t models.Table
	err = database.DB.QueryRow("SELECT id, number, capacity FROM tables WHERE id = ?", id).Scan(&t.ID, &t.Number, &t.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Table not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, t)
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
	var t models.Table
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("INSERT INTO tables (number, capacity) VALUES (?, ?)", t.Number, t.Capacity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, t)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	var t models.Table
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = database.DB.Exec("UPDATE tables SET number = ?, capacity = ? WHERE id = ?", t.Number, t.Capacity, id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, t)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	_, err = database.DB.Exec("DELETE FROM tables WHERE id = ?", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Table deleted successfully"})
}
