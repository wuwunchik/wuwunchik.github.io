package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/middleware"
	"wuwunchik.github.io/api/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Проверка учетных данных (заглушка, замените на реальную логику)
	if creds.Username != "admin" || creds.Password != "password" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := middleware.GenerateJWT(creds.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// api/controllers/auth.go
func Register(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	// Создаём пользователя с ролью "user" по умолчанию
	tx, err := database.DB.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer tx.Rollback()

	// 1. Добавляем пользователя
	result, err := tx.Exec(`
			INSERT INTO users (username, password_hash) VALUES (?, ?)
	`, user.Username, hashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// 2. Получаем ID нового пользователя
	userID, err := result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// 3. Назначаем роль "user" по умолчанию
	_, err = tx.Exec(`
			INSERT INTO user_roles (user_id, role_id)
			VALUES (?, (SELECT id FROM roles WHERE name = 'user'))
	`, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not assign role")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}
