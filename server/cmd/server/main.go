package main

import (
	"encoding/json"
	"log"
	"net/http"

	"vethub-go/internal/api"
	"vethub-go/internal/db"
	"vethub-go/internal/middleware"
)

func main() {
	database, err := db.Open()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/public/actuator/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
	})

	api.RegisterOwnerRoutes(mux, database)
	api.RegisterPetRoutes(mux, database)
	api.RegisterVisitRoutes(mux, database)
	api.RegisterVetRoutes(mux, database)

	handler := middleware.BasicAuth("user", "password", mux)
	handler = middleware.CORS("http://localhost:*", handler)

	log.Println("Starting VetHub Go server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
