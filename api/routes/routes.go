package routes

import (
	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
)

func RegisterRoutes(router *mux.Router) {
	// Маршруты для продуктов
	router.HandleFunc("/api/products/all", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/get/{id}", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/api/products/add", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products/update/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/delete/{id}", controllers.DeleteProduct).Methods("DELETE")

	// Маршруты для единиц измерения
	router.HandleFunc("/api/units/all", controllers.GetUnits).Methods("GET")
	router.HandleFunc("/api/units/get/{id}", controllers.GetUnit).Methods("GET")
	router.HandleFunc("/api/units/add", controllers.CreateUnit).Methods("POST")
	router.HandleFunc("/api/units/update/{id}", controllers.UpdateUnit).Methods("PUT")
	router.HandleFunc("/api/units/delete/{id}", controllers.DeleteUnit).Methods("DELETE")
}
