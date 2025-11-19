// backend/internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HostawayAPIKey string
	HostawayAccountID string
	HostawayBaseURL string
	ServerPort string
	UseMockData bool
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		HostawayAPIKey: getEnv("HOSTAWAY_API_KEY", "f94377ebbbb479490bb3ec364649168dc443dda2e4830facaf5de2e74ccc9152"),
		HostawayAccountID: getEnv("HOSTAWAY_ACCOUNT_ID", "61148"),
		HostawayBaseURL: getEnv("HOSTAWAY_BASE_URL", "https://api.hostaway.com/v1"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		UseMockData: getEnv("USE_MOCK_DATA", "false") == "true",
	}

	log.Printf("Configuration loaded: AccountID=%s, Port=%s, MockData=%v", 
		cfg.HostawayAccountID, cfg.ServerPort, cfg.UseMockData)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}