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

// GetUnits возвращает список всех единиц измерения
func GetUnits(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, abbreviation FROM units")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var units []models.Unit
	for rows.Next() {
		var u models.Unit
		err := rows.Scan(&u.ID, &u.Name, &u.Abbreviation)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		units = append(units, u)
	}

	utils.RespondWithJSON(w, http.StatusOK, units)
}

// GetUnit возвращает информацию о единице измерения
func GetUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	var u models.Unit
	err = database.DB.QueryRow("SELECT id, name, abbreviation FROM units WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Abbreviation)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Unit not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, u)
}

// CreateUnit создает новую единицу измерения
func CreateUnit(w http.ResponseWriter, r *http.Request) {
	var u models.Unit
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Начинаем транзакцию
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// Вставляем новую единицу измерения
	result, err := tx.Exec("INSERT INTO units (name, abbreviation) VALUES (?, ?)", u.Name, u.Abbreviation)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Получаем ID созданной единицы измерения
	id, _ := result.LastInsertId()
	u.ID = int(id)

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, u)
}

// UpdateUnit обновляет информацию о единице измерения
func UpdateUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	var u models.Unit
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Начинаем транзакцию
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// Обновляем единицу измерения
	_, err = tx.Exec("UPDATE units SET name = ?, abbreviation = ? WHERE id = ?", u.Name, u.Abbreviation, id)
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

	u.ID = id
	utils.RespondWithJSON(w, http.StatusOK, u)
}

// DeleteUnit удаляет единицу измерения
func DeleteUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid unit ID")
		return
	}

	// Начинаем транзакцию
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	// Удаляем единицу измерения
	_, err = tx.Exec("DELETE FROM units WHERE id = ?", id)
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

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Unit deleted successfully"})
}
