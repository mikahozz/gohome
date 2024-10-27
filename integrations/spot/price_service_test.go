package spot

import (
	"testing"
	"time"

	"github.com/mikahozz/gohome/integrations/spot/mock"
)

func TestGetSpotPrices(t *testing.T) {
	mockClient := mock.NewMockHTTPClient("testdata/oneDay.xml")
	spotService := NewSpotService(mockClient, "http://mock.api")

	periodStart, _ := time.Parse(time.RFC3339, "2024-10-22T21:00:00Z")
	periodEnd, _ := time.Parse(time.RFC3339, "2024-10-23T21:00:00Z")

	document, err := spotService.GetSpotPrices(periodStart, periodEnd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(document.TimeSeries) != 2 {
		t.Errorf("Expected 2 TimeSeries, got %d", len(document.TimeSeries))
	}
}

func TestGetSpotPrices_NoData(t *testing.T) {
	mockClient := mock.NewMockHTTPClient("testdata/noData_200.xml")
	spotService := NewSpotService(mockClient, "http://mock.api")

	periodStart := time.Now().AddDate(0, 0, 2)
	periodEnd := periodStart.AddDate(0, 0, 1)

	_, err := spotService.GetSpotPrices(periodStart, periodEnd)

	noDataErr, ok := err.(*NoDataError)
	if !ok {
		t.Fatalf("Expected NoDataError, got %T", err)
	}

	if noDataErr.Code != "999" {
		t.Errorf("Expected error code 999, got %s", noDataErr.Code)
	}

	expectedTextPrefix := "No matching data found for Data item ENERGY_PRICES and interval"
	if len(noDataErr.Text) < len(expectedTextPrefix) || noDataErr.Text[:len(expectedTextPrefix)] != expectedTextPrefix {
		t.Errorf("Expected error text to start with '%s', got '%s'", expectedTextPrefix, noDataErr.Text)
	}
}
