package main

import (
	"log"
	"net/http"

	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация базы данных
	database.InitDB()
	defer database.DB.Close()

	// Инициализация маршрутизатора
	router := mux.NewRouter()
	routes.RegisterRoutes(router, database.DB)

	// Запуск сервера
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
