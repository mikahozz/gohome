package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	apiEndpoint = "https://web-api.tp.entsoe.eu/api"
)

func getSpotPrices(periodStart, periodEnd time.Time) (*PublicationMarketDocument, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY not set in environment")
	}

	apiURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid API endpoint: %w", err)
	}

	params := url.Values{}
	params.Add("securityToken", apiKey)
	params.Add("documentType", "A44")
	params.Add("in_Domain", "10YFI-1--------U")
	params.Add("out_Domain", "10YFI-1--------U")
	params.Add("periodStart", periodStart.Format("200601021504"))
	params.Add("periodEnd", periodEnd.Format("200601021504"))

	apiURL.RawQuery = params.Encode()
	fmt.Println("Requesting url:", apiURL.String())

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading API response: %w", err)
	}

	var document PublicationMarketDocument
	err = xml.Unmarshal(body, &document)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling API response: %w", err)
	}

	return &document, nil
}

func parseISO8601Duration(duration string) (time.Duration, error) {
	// Remove the "PT" prefix
	duration = strings.TrimPrefix(duration, "PT")

	// Check if the duration ends with "M" for minutes
	if strings.HasSuffix(duration, "M") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(duration, "M"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration format: %s", duration)
		}
		return time.Duration(minutes) * time.Minute, nil
	}

	// Add more cases here if you need to handle other duration formats (e.g., hours, seconds)

	return 0, fmt.Errorf("unsupported duration format: %s", duration)
}

func convertToSpotPriceList(doc *PublicationMarketDocument, periodStart, periodEnd time.Time) (*SpotPriceList, error) {
	var spotPrices []SpotPrice

	for _, ts := range doc.TimeSeries {
		start, err := time.Parse("2006-01-02T15:04Z", ts.Period.TimeInterval.Start)
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %w", err)
		}

		resolution, err := parseISO8601Duration(ts.Period.Resolution)
		if err != nil {
			return nil, fmt.Errorf("error parsing resolution: %w", err)
		}

		for _, point := range ts.Period.Points {
			dateTime := start.Add(time.Duration(point.Position-1) * resolution)
			if dateTime.Before(periodStart) || dateTime.After(periodEnd) {
				continue
			}
			spotPrices = append(spotPrices, SpotPrice{
				DateTime: dateTime,
				Price:    point.Price,
			})
		}
	}

	// Sort the prices by datetime
	sort.Slice(spotPrices, func(i, j int) bool {
		return spotPrices[i].DateTime.Before(spotPrices[j].DateTime)
	})

	return &SpotPriceList{Prices: spotPrices}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	periodStart := time.Now().AddDate(0, 0, -1) // Yesterday
	periodEnd := time.Now()                     // Today

	document, err := getSpotPrices(periodStart, periodEnd)
	if err != nil {
		fmt.Println("Error retrieving spot prices:", err)
		return
	}

	spotPriceList, err := convertToSpotPriceList(document, periodStart, periodEnd)
	if err != nil {
		fmt.Println("Error converting spot prices:", err)
		return
	}

	jsonData, err := json.MarshalIndent(spotPriceList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
