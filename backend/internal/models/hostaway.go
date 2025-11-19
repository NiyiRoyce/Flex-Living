// backend/internal/models/hostaway.go
package models

// HostawayResponse represents the API response from Hostaway
type HostawayResponse struct {
	Status string           `json:"status"`
	Result []HostawayReview `json:"result"`
}

// HostawayReview represents a review from Hostaway API
type HostawayReview struct {
	ID             int              `json:"id"`
	Type           string           `json:"type"`
	Status         string           `json:"status"`
	Rating         *float64         `json:"rating"`
	PublicReview   string           `json:"publicReview"`
	ReviewCategory []ReviewCategory `json:"reviewCategory"`
	SubmittedAt    string           `json:"submittedAt"`
	GuestName      string           `json:"guestName"`
	ListingName    string           `json:"listingName"`
}

// ReviewCategory represents rating categories within a review
type ReviewCategory struct {
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
}