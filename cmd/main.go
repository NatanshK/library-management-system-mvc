package main

import (
	"fmt"
	"log"

	"library-management-system-mvc/config"
	"library-management-system-mvc/models"
)

func main() {
	fmt.Println("Booting up Library Management System...")

	config.ConnectDB()

	newUser := models.User{
		Username:      "student_01",
		Password:      "SuperSecretPassword123!",
		Email:         "student@university.edu",
		Role:          "client",
		RequestStatus: "not_requested",
	}

	fmt.Println("\n--- Testing Data Insertion ---")
	fmt.Printf("Attempting to register user: %s\n", newUser.Email)

	err := models.CreateUser(&newUser)
	if err != nil {

		log.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Println("User successfully created and secured in MySQL!")
	}

	fmt.Println("\n--- Testing Data Retrieval ---")

	fetchedUser, err := models.GetUserByEmail("student@university.edu")
	if err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}

	fmt.Println("User found in database!")
	fmt.Printf("ID: %d\n", fetchedUser.ID)
	fmt.Printf("Username: %s\n", fetchedUser.Username)
	fmt.Printf("Role: %s\n", fetchedUser.Role)
	fmt.Printf("Stored Hash: %s...\n", fetchedUser.Password[:15]) // Print just the first 15 chars
	fmt.Printf("Stored Salt: %s\n", fetchedUser.Salt)
}
