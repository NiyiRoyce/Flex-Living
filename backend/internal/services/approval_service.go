// backend/internal/services/approval_service.go
package services

import (
	"Flex-Living/internal/models"
	"fmt"
	"sync"
	"time"
)

type ApprovalService struct {
	reviews map[string]*models.NormalizedReview
	mu      sync.RWMutex
}

func NewApprovalService() *ApprovalService {
	return &ApprovalService{
		reviews: make(map[string]*models.NormalizedReview),
	}
}

// LoadReviews initializes the service with reviews
func (s *ApprovalService) LoadReviews(reviews []models.NormalizedReview) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range reviews {
		s.reviews[reviews[i].ID] = &reviews[i]
	}
}

// ApproveReview marks a review as approved
func (s *ApprovalService) ApproveReview(reviewID, approvedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, exists := s.reviews[reviewID]
	if !exists {
		return fmt.Errorf("review not found: %s", reviewID)
	}

	now := time.Now()
	review.Status = "approved"
	review.ApprovalStatus.IsApproved = true
	review.ApprovalStatus.IsRejected = false
	review.ApprovalStatus.ApprovedAt = &now
	review.ApprovalStatus.ApprovedBy = approvedBy

	return nil
}

// RejectReview marks a review as rejected
func (s *ApprovalService) RejectReview(reviewID, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, exists := s.reviews[reviewID]
	if !exists {
		return fmt.Errorf("review not found: %s", reviewID)
	}

	now := time.Now()
	review.Status = "rejected"
	review.ApprovalStatus.IsApproved = false
	review.ApprovalStatus.IsRejected = true
	review.ApprovalStatus.RejectedAt = &now
	review.ApprovalStatus.RejectionReason = reason

	return nil
}

// GetReview retrieves a specific review
func (s *ApprovalService) GetReview(reviewID string) (*models.NormalizedReview, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	review, exists := s.reviews[reviewID]
	if !exists {
		return nil, fmt.Errorf("review not found: %s", reviewID)
	}

	return review, nil
}

// GetAllReviews returns all reviews
func (s *ApprovalService) GetAllReviews() []models.NormalizedReview {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reviews := make([]models.NormalizedReview, 0, len(s.reviews))
	for _, review := range s.reviews {
		reviews = append(reviews, *review)
	}

	return reviews
}

// GetReviewStats calculates aggregate statistics
func (s *ApprovalService) GetReviewStats() models.ReviewStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := models.ReviewStats{
		RatingsBySource:  make(map[string]float64),
		CategoryAverages: make(map[string]float64),
		StatusBreakdown:  make(map[string]int),
	}

	if len(s.reviews) == 0 {
		return stats
	}

	var totalRating float64
	sourceRatings := make(map[string][]float64)
	categoryRatings := make(map[string][]float64)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, review := range s.reviews {
		stats.TotalReviews++
		totalRating += review.Rating
		
		// Source ratings
		sourceRatings[review.Source] = append(sourceRatings[review.Source], review.Rating)
		
		// Category ratings
		for category, rating := range review.Categories {
			categoryRatings[category] = append(categoryRatings[category], rating)
		}
		
		// Status breakdown
		stats.StatusBreakdown[review.Status]++
		
		// Recent reviews
		if review.SubmittedAt.After(thirtyDaysAgo) {
			stats.RecentReviews++
		}
	}

	// Calculate averages
	stats.AverageRating = totalRating / float64(stats.TotalReviews)

	for source, ratings := range sourceRatings {
		var sum float64
		for _, rating := range ratings {
			sum += rating
		}
		stats.RatingsBySource[source] = sum / float64(len(ratings))
	}

	for category, ratings := range categoryRatings {
		var sum float64
		for _, rating := range ratings {
			sum += rating
		}
		stats.CategoryAverages[category] = sum / float64(len(ratings))
	}

	return stats
}