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

func GetPrices(start, end time.Time) (*SpotPriceList, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not set in environment")
	}

	client := NewDefaultHTTPClient(apiKey)
	spotService := NewSpotService(client, apiEndpoint)

	document, err := spotService.GetSpotPrices(start, end)
	if err != nil {
		return nil, fmt.Errorf("error getting spot prices: %w", err)
	}

	return ConvertToSpotPriceList(document, start, end)
}
