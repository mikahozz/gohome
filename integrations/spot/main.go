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

func getPrices() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}

	client := NewDefaultHTTPClient(apiKey)
	spotService := NewSpotService(client, apiEndpoint)

	prices, err := spotService.GetSpotPrices(time.Now(), time.Now().Add(24*time.Hour))
	if err != nil {
		log.Fatalf("Error getting spot prices: %v", err)
	}
	fmt.Println(prices)

	select {} // Block forever
}
