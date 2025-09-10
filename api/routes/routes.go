package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/controllers"
)

func RegisterRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/api/products", controllers.GetProducts(db)).Methods("GET")
	router.HandleFunc("/api/products/{id}", controllers.GetProduct(db)).Methods("GET")
}
