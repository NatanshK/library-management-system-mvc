package main

import (
	"fmt"
	"log"

	"library-management-system-mvc/config"
	"library-management-system-mvc/models" // Import your new models package
)

func main() {
	fmt.Println("Booting up Library Management System...")

	// 1. Fire up the database connection pool
	config.ConnectDB()

	// 2. Create a dummy user blueprint in memory
	newUser := models.User{
		Username:      "student_01",
		Password:      "SuperSecretPassword123!",
		Email:         "student@university.edu",
		Role:          "client",
		RequestStatus: "not_requested",
	}

	fmt.Println("\n--- Testing Data Insertion ---")
	fmt.Printf("Attempting to register user: %s\n", newUser.Email)

	// 3. Call the insertion model
	err := models.CreateUser(&newUser)
	if err != nil {
		// If the user already exists, this will print a MySQL duplicate key error
		log.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Println("✅ User successfully created and secured in MySQL!")
	}

	fmt.Println("\n--- Testing Data Retrieval ---")

	// 4. Call the read model to fetch the user back out
	fetchedUser, err := models.GetUserByEmail("student@university.edu")
	if err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}

	// 5. Print the retrieved data to the terminal
	fmt.Println("✅ User found in database!")
	fmt.Printf("ID: %d\n", fetchedUser.ID)
	fmt.Printf("Username: %s\n", fetchedUser.Username)
	fmt.Printf("Role: %s\n", fetchedUser.Role)
	fmt.Printf("Stored Hash: %s...\n", fetchedUser.Password[:15]) // Print just the first 15 chars
	fmt.Printf("Stored Salt: %s\n", fetchedUser.Salt)
}
