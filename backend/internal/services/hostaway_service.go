// backend/internal/services/hostaway_service.go
package services

import (
	"encoding/json"
	"Flex-Living/internal/config"
	"Flex-Living/internal/models"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HostawayService struct {
	config *config.Config
	client *http.Client
}

func NewHostawayService(cfg *config.Config) *HostawayService {
	return &HostawayService{
		config: cfg,
		client: &http.Client{},
	}
}

// FetchReviews retrieves reviews from Hostaway API or mock data
func (s *HostawayService) FetchReviews() ([]models.HostawayReview, error) {
	// Use mock data if configured
	if s.config.UseMockData {
		return s.fetchMockReviews()
	}

	// Construct API URL
	url := fmt.Sprintf("%s/reviews?accountId=%s", 
		s.config.HostawayBaseURL, 
		s.config.HostawayAccountID)

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Add("Authorization", "Bearer "+s.config.HostawayAPIKey)
	req.Header.Add("Content-Type", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResponse models.HostawayResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResponse.Result, nil
}

// fetchMockReviews loads reviews from mock data file
func (s *HostawayService) fetchMockReviews() ([]models.HostawayReview, error) {
	file, err := os.Open("mockdata/hostaway_mock.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open mock data: %w", err)
	}
	defer file.Close()

	var apiResponse models.HostawayResponse
	if err := json.NewDecoder(file).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode mock data: %w", err)
	}

	return apiResponse.Result, nil
}