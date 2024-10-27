package spot

import (
	"encoding/xml"
	"os"
	"testing"
	"time"
)

func TestConvertToSpotPriceList(t *testing.T) {
	// Read the test XML file
	xmlData, err := os.ReadFile("mock/oneDay.xml")
	if err != nil {
		t.Fatalf("Failed to read test XML file: %v", err)
	}

	// Unmarshal the XML data
	var doc PublicationMarketDocument
	err = xml.Unmarshal(xmlData, &doc)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML data: %v", err)
	}

	// Set the period start and end times
	periodStart, _ := time.Parse(time.RFC3339, "2024-10-22T21:00:00Z")
	periodEnd, _ := time.Parse(time.RFC3339, "2024-10-23T21:00:00Z")

	// Convert to SpotPriceList
	spotPriceList, err := ConvertToSpotPriceList(&doc, periodStart, periodEnd, time.UTC)
	if err != nil {
		t.Fatalf("Failed to convert to SpotPriceList: %v", err)
	}

	// Check the number of prices
	expectedCount := 24 // 24 hours
	if len(spotPriceList.Prices) != expectedCount {
		t.Errorf("Expected %d prices, but got %d", expectedCount, len(spotPriceList.Prices))
	}

	// Check the first price
	expectedFirstPrice := SpotPrice{
		DateTime:  time.Date(2024, 10, 22, 21, 0, 0, 0, time.UTC),
		PriceCkwh: -0.08,
	}
	if spotPriceList.Prices[0] != expectedFirstPrice {
		t.Errorf("Expected first price %+v, but got %+v", expectedFirstPrice, spotPriceList.Prices[0])
	}

	// Check the last price
	expectedLastPrice := SpotPrice{
		DateTime:  time.Date(2024, 10, 23, 21, 0, 0, 0, time.UTC),
		PriceCkwh: -0.081,
	}
	if spotPriceList.Prices[len(spotPriceList.Prices)-1] != expectedLastPrice {
		t.Errorf("Expected last price %+v, but got %+v", expectedLastPrice, spotPriceList.Prices[len(spotPriceList.Prices)-1])
	}
}
