package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"wuwunchik.github.io/api/database"
	"wuwunchik.github.io/api/routes"
)

func main() {
	database.InitDB()
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
