// backend/internal/handlers/reviews_handler.go
package handlers

import (
	"encoding/json"
	"Flex-Living/internal/config"
	"Flex-Living/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ReviewsHandler struct {
	hostawayService      *services.HostawayService
	normalizationService *services.NormalizationService
	approvalService      *services.ApprovalService
}

func NewReviewsHandler(cfg *config.Config) *ReviewsHandler {
	handler := &ReviewsHandler{
		hostawayService:      services.NewHostawayService(cfg),
		normalizationService: services.NewNormalizationService(),
		approvalService:      services.NewApprovalService(),
	}

	// Initialize with reviews
	handler.initializeReviews()

	return handler
}

// initializeReviews loads and normalizes reviews on startup
func (h *ReviewsHandler) initializeReviews() {
	reviews, err := h.hostawayService.FetchReviews()
	if err != nil {
		log.Printf("Warning: Failed to initialize reviews: %v", err)
		return
	}

	normalized := h.normalizationService.NormalizeHostawayReviews(reviews)
	h.approvalService.LoadReviews(normalized)
	
	log.Printf("Initialized with %d reviews", len(normalized))
}

// GetHostawayReviews returns raw Hostaway reviews
func (h *ReviewsHandler) GetHostawayReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.hostawayService.FetchReviews()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"count":  len(reviews),
		"data":   reviews,
	})
}

// GetNormalizedReviews returns normalized reviews with optional filtering
func (h *ReviewsHandler) GetNormalizedReviews(w http.ResponseWriter, r *http.Request) {
	reviews := h.approvalService.GetAllReviews()

	// Filter by status if provided
	status := r.URL.Query().Get("status")
	if status != "" {
		filtered := make([]interface{}, 0)
		for _, review := range reviews {
			if review.Status == status {
				filtered = append(filtered, review)
			}
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"count":  len(filtered),
			"data":   filtered,
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"count":  len(reviews),
		"data":   reviews,
	})
}

// ApproveReview handles review approval
func (h *ReviewsHandler) ApproveReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["id"]

	var req struct {
		ApprovedBy string `json:"approvedBy"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ApprovedBy == "" {
		req.ApprovedBy = "admin"
	}

	if err := h.approvalService.ApproveReview(reviewID, req.ApprovedBy); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	review, _ := h.approvalService.GetReview(reviewID)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Review approved successfully",
		"data":    review,
	})
}

// RejectReview handles review rejection
func (h *ReviewsHandler) RejectReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["id"]

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.approvalService.RejectReview(reviewID, req.Reason); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	review, _ := h.approvalService.GetReview(reviewID)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Review rejected successfully",
		"data":    review,
	})
}

// GetReviewStats returns aggregated statistics
func (h *ReviewsHandler) GetReviewStats(w http.ResponseWriter, r *http.Request) {
	stats := h.approvalService.GetReviewStats()

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   stats,
	})
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{
		"status": "error",
		"error":  message,
	})
}