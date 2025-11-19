// backend/internal/models/normalized.go
package models

import "time"

// NormalizedReview represents a standardized review format
type NormalizedReview struct {
	ID             string             `json:"id"`
	Source         string             `json:"source"` // "hostaway", "google", etc.
	PropertyID     string             `json:"propertyId"`
	PropertyName   string             `json:"propertyName"`
	GuestName      string             `json:"guestName"`
	Rating         float64            `json:"rating"`
	ReviewText     string             `json:"reviewText"`
	Categories     map[string]float64 `json:"categories"`
	SubmittedAt    time.Time          `json:"submittedAt"`
	Status         string             `json:"status"` // "pending", "approved", "rejected"
	ApprovalStatus ApprovalStatus     `json:"approvalStatus"`
}

// ApprovalStatus tracks review approval workflow
type ApprovalStatus struct {
	IsApproved   bool       `json:"isApproved"`
	IsRejected   bool       `json:"isRejected"`
	ApprovedAt   *time.Time `json:"approvedAt,omitempty"`
	RejectedAt   *time.Time `json:"rejectedAt,omitempty"`
	ApprovedBy   string     `json:"approvedBy,omitempty"`
	RejectionReason string  `json:"rejectionReason,omitempty"`
}

// ReviewStats provides aggregated review statistics
type ReviewStats struct {
	TotalReviews    int                `json:"totalReviews"`
	AverageRating   float64            `json:"averageRating"`
	RatingsBySource map[string]float64 `json:"ratingsBySource"`
	CategoryAverages map[string]float64 `json:"categoryAverages"`
	StatusBreakdown map[string]int     `json:"statusBreakdown"`
	RecentReviews   int                `json:"recentReviews"` // Last 30 days
}