package sun

import (
	"os"
	"path/filepath"
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
	// Create test data with different years to verify year-independent behavior
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
			name:           "Single day - same year",
			startDate:      "2025-01-01",
			endDate:        nil,
			expectedDates:  []string{"2025-01-01"},
			expectedLength: 1,
		},
		{
			name:           "Single day - different year",
			startDate:      "2023-01-01", // Different year, same month/day
			endDate:        nil,
			expectedDates:  []string{"2025-01-01"},
			expectedLength: 1,
		},
		{
			name:           "Date range - same year",
			startDate:      "2025-01-02",
			endDate:        stringPtr("2025-01-03"),
			expectedDates:  []string{"2025-01-02", "2025-01-03"},
			expectedLength: 2,
		},
		{
			name:           "Date range - different years",
			startDate:      "2023-01-02",            // Different year, same month/day
			endDate:        stringPtr("2024-01-03"), // Different year, same month/day
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

func TestGetSunriseAndSunsetToday(t *testing.T) {
	// Change working directory to repo root so relative path in GetSunrise/Set functions resolves
	wd, _ := os.Getwd()
	// Assume test file located at .../integrations/sun; repo root is two levels up
	root := filepath.Clean(filepath.Join(wd, "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("failed to chdir root: %v", err)
	}
	defer os.Chdir(wd)

	sunrise := GetSunriseToday()
	sunset := GetSunsetToday()
	if !sunset.After(sunrise) {
		t.Fatalf("expected sunset (%v) to be after sunrise (%v)", sunset, sunrise)
	}
	// Load raw data to verify the parsed times match the expected date
	sunData, err := LoadSunData("integrations/sun/sun_helsinki_2025.json")
	if err != nil {
		t.Fatalf("could not load sun data: %v", err)
	}
	todayArr := sunData.GetDailyData(time.Now(), nil)
	if len(todayArr) == 0 {
		t.Fatalf("no data for today in test dataset")
	}
	d := todayArr[0]
	// Ensure the date portion matches
	if sunrise.Format("2006-01-02") != d.Date || sunset.Format("2006-01-02") != d.Date {
		t.Fatalf("parsed times date mismatch. got sunrise date=%s sunset date=%s want=%s", sunrise.Format("2006-01-02"), sunset.Format("2006-01-02"), d.Date)
	}
}

// Helper function to get pointer to string
func stringPtr(s string) *string {
	return &s
}
