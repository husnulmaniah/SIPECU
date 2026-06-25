package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sipecut/config"
	"sipecut/routes"
	"sipecut/services"
	"sipecut/storage"
)

func main() {
	log.Println("Starting SIPECUT Local Development Server...")

	// 1. Connect to Database & Run Auto-Migrations
	config.ConnectDatabase()

	// 2. Initialize File Storage Configuration
	storage.InitStorage()

	// 3. Run Initial Retirement Check on active employees
	db := config.GetDB()
	if db != nil {
		log.Println("Running initial employee retirement check...")
		if err := services.CheckAllEmployeesRetirementStatus(db); err != nil {
			log.Printf("Error during retirement check: %v\n", err)
		}
	}

	// 4. Start simple hourly background routine to check retirement status
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			dbCon := config.GetDB()
			if dbCon != nil {
				log.Println("Scheduled background check: evaluating employee retirement status...")
				_ = services.CheckAllEmployeesRetirementStatus(dbCon)
			}
		}
	}()

	// 5. Setup Gin Router
	r := routes.SetupRouter()

	// 6. Bind to Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SIPECUT Backend listening on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}
}
