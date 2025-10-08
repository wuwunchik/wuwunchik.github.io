// api/middleware/roles.go
package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
	"wuwunchik.github.io/api/middleware"
	"wuwunchik.github.io/api/utils"
)

// RoleCheck проверяет, есть ли у пользователя хотя бы одна из требуемых ролей
func RoleCheck(requiredRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*Claims)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			hasRole := false
			for _, requiredRole := range requiredRoles {
				for _, role := range claims.Roles {
					if role == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				utils.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// api/routes/routes.go
func RegisterRoutes(router *mux.Router) {
	// Маршруты аутентификации
	router.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", controllers.Register).Methods("POST")

	// Публичные маршруты (без JWT)
	router.HandleFunc("/api/products/public", controllers.GetPublicProducts).Methods("GET")

	// Маршруты для авторизованных пользователей
	router.HandleFunc("/api/products/all",
		middleware.ValidateJWT(controllers.GetProducts)).Methods("GET")

	// Маршруты только для админов
	// Маршруты только для админов
	router.HandleFunc(
		"/api/users/all",
		middleware.ValidateJWT(
			middleware.RoleCheck("admin")(controllers.GetAllUsers),
		),
	).Methods("GET")

	// Маршруты для админов или менеджеров
	// Маршруты для админов или менеджеров
	router.HandleFunc(
		"/api/products/add",
		middleware.ValidateJWT(
			middleware.RoleCheck("admin", "manager")(controllers.CreateProduct),
		),
	).Methods("POST")

}
