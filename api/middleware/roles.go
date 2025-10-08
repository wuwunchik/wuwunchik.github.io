// api/middleware/roles.go
package middleware

import (
	"net/http"

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
