// backend/cmd/server/main.go
package main

import (
	"Flex-Living/internal/config"
	"Flex-Living/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize router
	r := mux.NewRouter()

	// Initialize handlers
	reviewsHandler := handlers.NewReviewsHandler(cfg)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/reviews/hostaway", reviewsHandler.GetHostawayReviews).Methods("GET")
	api.HandleFunc("/reviews/normalized", reviewsHandler.GetNormalizedReviews).Methods("GET")
	api.HandleFunc("/reviews/{id}/approve", reviewsHandler.ApproveReview).Methods("POST")
	api.HandleFunc("/reviews/{id}/reject", reviewsHandler.RejectReview).Methods("POST")
	api.HandleFunc("/reviews/stats", reviewsHandler.GetReviewStats).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// Start server
	port := cfg.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}