package main

import (
	"log"
	"scandiff/api"
	"scandiff/db"
)

func main() {
	// Initialize database
	database, err := db.InitDB("./data/snapshots.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize and run the API server
	router := api.SetupRouter(database)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
