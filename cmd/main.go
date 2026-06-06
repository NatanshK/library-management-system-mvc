package main

import (
	"fmt"
	"log"
	"net/http"

	"library-management-system-mvc/config"
	"library-management-system-mvc/controllers"
)

func main() {
	// 1. Connect to the database
	config.ConnectDB()

	// 2. Auth Routes
	http.HandleFunc("/api/register", controllers.Register)
	http.HandleFunc("/api/login", controllers.Login)

	// 3. Book Routes (Currently Unprotected)
	// 4. Transaction Routes
	http.HandleFunc("/api/transactions/request", controllers.RequestCheckout)
	http.HandleFunc("/api/transactions/approve", controllers.ApproveCheckout)
	http.HandleFunc("/api/transactions/history", controllers.GetUserHistory)
	http.HandleFunc("/api/transactions/queue", controllers.GetAdminQueue)

	// 4. Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
