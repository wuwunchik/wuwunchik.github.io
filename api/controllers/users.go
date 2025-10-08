package controllers

import (
	"net/http"

	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/utils"
)

// GetAllUsers возвращает список всех пользователей с их ролями (только для админов)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Роль "admin" уже проверена в middleware.RoleCheck

	rows, err := database.DB.Query(`
			SELECT u.id, u.username, u.created_at, GROUP_CONCAT(r.name, ', ') as roles
			FROM users u
			LEFT JOIN user_roles ur ON u.id = ur.user_id
			LEFT JOIN roles r ON ur.role_id = r.id
			GROUP BY u.id
	`)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	type UserWithRoles struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
		Roles     string `json:"roles"`
	}

	var users []UserWithRoles
	for rows.Next() {
		var u UserWithRoles
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt, &u.Roles); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		users = append(users, u)
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}
