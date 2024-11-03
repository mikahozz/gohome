package spot

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	apiEndpoint = "https://web-api.tp.entsoe.eu/api"
)

func GetPrices(start, end time.Time, location *time.Location) (*SpotPriceList, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not set in environment")
	}

	start = start.In(location)
	end = end.In(location)

	client := NewDefaultHTTPClient(apiKey)
	spotService := NewSpotService(client, apiEndpoint)

	return spotService.GetSpotPrices(start, end)
}
