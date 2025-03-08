package sun

import (
	"encoding/json"
	"os"
	"time"
)

// SunData represents the overall structure of the sun data JSON
type SunData struct {
	Results []DailyData `json:"results"`
}

// DailyData represents the sun data for a single day
type DailyData struct {
	Date       string `json:"date"`
	Sunrise    string `json:"sunrise"`
	Sunset     string `json:"sunset"`
	FirstLight string `json:"first_light"`
	LastLight  string `json:"last_light"`
	Dawn       string `json:"dawn"`
	Dusk       string `json:"dusk"`
	SolarNoon  string `json:"solar_noon"`
	GoldenHour string `json:"golden_hour"`
	DayLength  string `json:"day_length"`
	Timezone   string `json:"timezone"`
	UTCOffset  int    `json:"utc_offset"`
}

// LoadSunData loads sun data from a JSON file
func LoadSunData(filePath string) (*SunData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var sunData SunData
	if err := json.Unmarshal(data, &sunData); err != nil {
		return nil, err
	}

	return &sunData, nil
}

// GetDailyData returns the sun data for a specific date or date range
func (s *SunData) GetDailyData(startDate time.Time, endDate *time.Time) []*DailyData {
	startDateStr := startDate.Format("2006-01-02")
	var results []*DailyData

	// Determine end date string - either the provided end date or the start date if no end date
	endDateStr := startDateStr
	if endDate != nil {
		endDateStr = endDate.Format("2006-01-02")
	}

	// Collect all dates in the range
	for i, daily := range s.Results {
		if daily.Date >= startDateStr && daily.Date <= endDateStr {
			results = append(results, &s.Results[i])
		}
	}

	return results
}
