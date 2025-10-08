package routes

import (
	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
	"wuwunchik.github.io/api/middleware"
)

// func RegisterRoutes(router *mux.Router) {
// 	// Маршруты для продуктов
// 	router.HandleFunc("/api/products/all", controllers.GetProducts).Methods("GET")
// 	router.HandleFunc("/api/products/get/{id}", controllers.GetProduct).Methods("GET")
// 	router.HandleFunc("/api/products/add", controllers.CreateProduct).Methods("POST")
// 	router.HandleFunc("/api/products/update/{id}", controllers.UpdateProduct).Methods("PUT")
// 	router.HandleFunc("/api/products/delete/{id}", controllers.DeleteProduct).Methods("DELETE")

// 	// Маршруты для единиц измерения
// 	router.HandleFunc("/api/units/all", controllers.GetUnits).Methods("GET")
// 	router.HandleFunc("/api/units/get/{id}", controllers.GetUnit).Methods("GET")
// 	router.HandleFunc("/api/units/add", controllers.CreateUnit).Methods("POST")
// 	router.HandleFunc("/api/units/update/{id}", controllers.UpdateUnit).Methods("PUT")
// 	router.HandleFunc("/api/units/delete/{id}", controllers.DeleteUnit).Methods("DELETE")
// }

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
