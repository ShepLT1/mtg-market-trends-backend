package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"backend/api"
	"backend/data"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	// Connect to Postgres
	data.Connect()

	sqlDB, err := data.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("✅ Connected to database")

	cardRepo := data.NewCardRepository(data.DB)
	listingRepo := data.NewListingRepository(data.DB)

	// Initialize router (API endpoints)
	router := api.NewRouter(cardRepo, listingRepo)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
