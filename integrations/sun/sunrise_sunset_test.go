package sun

import (
	"testing"
	"time"
)

func TestLoadSunData(t *testing.T) {
	// Test loading the data
	sunData, err := LoadSunData("sun_helsinki_2025.json")
	if err != nil {
		t.Fatalf("LoadSunData failed: %v", err)
	}

	// Verify the data was loaded correctly
	if len(sunData.Results) != 365 {
		t.Errorf("Expected 365 result, got %d", len(sunData.Results))
	}

	result := sunData.Results[0]
	expected := DailyData{
		Date:      "2025-01-01",
		Sunrise:   "9:25:22 AM",
		Timezone:  "Europe/Helsinki",
		UTCOffset: 120,
	}
	if result.Date != expected.Date || result.Sunrise != expected.Sunrise ||
		result.Timezone != expected.Timezone || result.UTCOffset != expected.UTCOffset {
		t.Errorf("First result mismatch.\nGot: %+v\nWant: %+v", result, expected)
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
			{
				Date:    "2025-01-03",
				Sunrise: "9:24:15 AM",
				Sunset:  "3:27:30 PM",
			},
		},
	}

	// Test cases
	testCases := []struct {
		name           string
		startDate      string
		endDate        *string
		expectedDates  []string
		expectedLength int
	}{
		{
			name:           "Single day",
			startDate:      "2025-01-01",
			endDate:        nil,
			expectedDates:  []string{"2025-01-01"},
			expectedLength: 1,
		},
		{
			name:           "Date range",
			startDate:      "2025-01-02",
			endDate:        stringPtr("2025-01-03"),
			expectedDates:  []string{"2025-01-02", "2025-01-03"},
			expectedLength: 2,
		},
		{
			name:           "Non-existent date",
			startDate:      "2025-01-04",
			endDate:        nil,
			expectedDates:  []string{},
			expectedLength: 0,
		},
		{
			name:           "Partial range",
			startDate:      "2025-01-03",
			endDate:        stringPtr("2025-01-05"),
			expectedDates:  []string{"2025-01-03"},
			expectedLength: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			startDate, _ := time.Parse("2006-01-02", tc.startDate)
			var endDate *time.Time
			if tc.endDate != nil {
				parsed, _ := time.Parse("2006-01-02", *tc.endDate)
				endDate = &parsed
			}

			result := sunData.GetDailyData(startDate, endDate)

			// Check length
			if len(result) != tc.expectedLength {
				t.Fatalf("Expected to find %d results, got %d", tc.expectedLength, len(result))
			}

			// Check dates if we expect results
			for i, expectedDate := range tc.expectedDates {
				if i < len(result) && result[i].Date != expectedDate {
					t.Errorf("Expected date %s at position %d, got %s", expectedDate, i, result[i].Date)
				}
			}
		})
	}
}

// Helper function to get pointer to string
func stringPtr(s string) *string {
	return &s
}
