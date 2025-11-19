// backend/internal/services/normalization_service.go
package services

import (
	"Flex-Living/internal/models"
	"fmt"
	"time"
)

type NormalizationService struct{}

func NewNormalizationService() *NormalizationService {
	return &NormalizationService{}
}

// NormalizeHostawayReviews converts Hostaway reviews to normalized format
func (s *NormalizationService) NormalizeHostawayReviews(reviews []models.HostawayReview) []models.NormalizedReview {
	normalized := make([]models.NormalizedReview, 0, len(reviews))

	for _, review := range reviews {
		normalized = append(normalized, s.normalizeHostawayReview(review))
	}

	return normalized
}

// normalizeHostawayReview converts a single Hostaway review
func (s *NormalizationService) normalizeHostawayReview(review models.HostawayReview) models.NormalizedReview {
	// Calculate overall rating from categories
	overallRating := s.calculateOverallRating(review)

	// Parse submission date
	submittedAt, _ := time.Parse("2006-01-02 15:04:05", review.SubmittedAt)

	// Convert categories to map
	categories := make(map[string]float64)
	for _, cat := range review.ReviewCategory {
		categories[cat.Category] = cat.Rating
	}

	// Extract property ID from listing name (simplified)
	propertyID := s.extractPropertyID(review.ListingName)

	return models.NormalizedReview{
		ID:           fmt.Sprintf("hostaway-%d", review.ID),
		Source:       "hostaway",
		PropertyID:   propertyID,
		PropertyName: review.ListingName,
		GuestName:    review.GuestName,
		Rating:       overallRating,
		ReviewText:   review.PublicReview,
		Categories:   categories,
		SubmittedAt:  submittedAt,
		Status:       s.mapStatus(review.Status),
		ApprovalStatus: models.ApprovalStatus{
			IsApproved: review.Status == "published",
			IsRejected: review.Status == "rejected",
		},
	}
}

// calculateOverallRating computes average from category ratings
func (s *NormalizationService) calculateOverallRating(review models.HostawayReview) float64 {
	if review.Rating != nil {
		return *review.Rating
	}

	if len(review.ReviewCategory) == 0 {
		return 0
	}

	var sum float64
	for _, cat := range review.ReviewCategory {
		sum += cat.Rating
	}

	return sum / float64(len(review.ReviewCategory))
}

// extractPropertyID extracts a property identifier from listing name
func (s *NormalizationService) extractPropertyID(listingName string) string {
	// Simple extraction - in production, use regex or database mapping
	if listingName == "" {
		return "unknown"
	}
	return fmt.Sprintf("prop-%d", len(listingName)) // Simplified ID generation
}

// mapStatus converts Hostaway status to normalized status
func (s *NormalizationService) mapStatus(hostawayStatus string) string {
	switch hostawayStatus {
	case "published":
		return "approved"
	case "rejected":
		return "rejected"
	default:
		return "pending"
	}
}