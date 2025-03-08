package sun

import (
	"testing"
	"time"
)

func TestLoadSunData(t *testing.T) {
	// Test loading the data
	sunData, err := LoadSunData("sun_helsinki.json")
	if err != nil {
		t.Fatalf("LoadSunData failed: %v", err)
	}

	// Verify the data was loaded correctly
	if len(sunData.Results) != 365 {
		t.Errorf("Expected 365 result, got %d", len(sunData.Results))
	}

	result := sunData.Results[0]
	if result.Date != "2025-01-01" {
		t.Errorf("Expected date 2025-01-01, got %s", result.Date)
	}
	if result.Sunrise != "9:25:22 AM" {
		t.Errorf("Expected sunrise 9:25:22 AM, got %s", result.Sunrise)
	}
	if result.Timezone != "Europe/Helsinki" {
		t.Errorf("Expected timezone Europe/Helsinki, got %s", result.Timezone)
	}
	if result.UTCOffset != 120 {
		t.Errorf("Expected UTC offset 120, got %d", result.UTCOffset)
	}
}

func TestGetDailyData(t *testing.T) {
	// Create test data
	sunData := &SunData{
		Results: []DailyData{
			{
				Date:    "2025-01-01",
				Sunrise: "9:25:22 AM",
				Sunset:  "3:24:31 PM",
			},
			{
				Date:    "2025-01-02",
				Sunrise: "9:24:50 AM",
				Sunset:  "3:25:58 PM",
			},
		},
	}

	// Test finding existing date
	date1, _ := time.Parse("2006-01-02", "2025-01-01")
	result := sunData.GetDailyData(date1)
	if result == nil {
		t.Fatal("Expected to find data for 2025-01-01, got nil")
	}
	if result.Date != "2025-01-01" {
		t.Errorf("Expected date 2025-01-01, got %s", result.Date)
	}

	// Test finding another existing date
	date2, _ := time.Parse("2006-01-02", "2025-01-02")
	result = sunData.GetDailyData(date2)
	if result == nil {
		t.Fatal("Expected to find data for 2025-01-02, got nil")
	}
	if result.Date != "2025-01-02" {
		t.Errorf("Expected date 2025-01-02, got %s", result.Date)
	}

	// Test with non-existent date
	date3, _ := time.Parse("2006-01-02", "2025-01-03")
	result = sunData.GetDailyData(date3)
	if result != nil {
		t.Errorf("Expected nil for non-existent date, got %+v", result)
	}
}
