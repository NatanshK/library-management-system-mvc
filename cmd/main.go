package main

import (
	"log"
	"net/http"

	"library-management-system-mvc/config"
	"library-management-system-mvc/routes"
)

func main() {

	config.ConnectDB()

	router := routes.SetupRoutes()

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
