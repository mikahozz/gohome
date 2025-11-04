package sun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
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
// The year part of the date is ignored, only month and day are considered
func (s *SunData) GetDailyData(startDate time.Time, endDate *time.Time) []*DailyData {
	// Format dates as MM-DD to ignore year
	startMonthDay := startDate.Format("01-02")
	var results []*DailyData

	// Determine end date - either the provided end date or the start date if no end date
	endMonthDay := startMonthDay
	if endDate != nil {
		endMonthDay = endDate.Format("01-02")
	}

	// Collect all dates in the range
	for i, daily := range s.Results {
		// Parse the date from the data to extract month and day
		t, err := time.Parse("2006-01-02", daily.Date)
		if err != nil {
			continue // Skip invalid dates
		}

		// Compare only month and day
		dailyMonthDay := t.Format("01-02")

		if dailyMonthDay >= startMonthDay && dailyMonthDay <= endMonthDay {
			results = append(results, &s.Results[i])
		}
	}

	return results
}

// getTodayTime is an internal helper to fetch and parse a specific field (sunrise/sunset)
// for today's date from the static JSON file. It panics on error to preserve existing
// public function behaviour.
func getTodayTime(field string) time.Time {
	sunData, err := LoadSunData("integrations/sun/sun_helsinki_2025.json")
	if err != nil {
		log.Error().Err(err).Msg("Error loading sun data")
		panic(err)
	}
	sArr := sunData.GetDailyData(time.Now(), nil)
	if len(sArr) == 0 {
		e := errors.New("no sun data found for the date")
		log.Error().Err(e).Msg("No sun data found for the date")
		panic(e)
	}
	s := sArr[0]
	var raw string
	switch field {
	case "sunrise":
		raw = s.Sunrise
	case "sunset":
		raw = s.Sunset
	default:
		e := errors.New("unsupported field: " + field)
		panic(e)
	}
	tm, err := time.Parse("2006-01-02 3:04:05 PM", fmt.Sprintf("%s %s", s.Date, raw))
	if err != nil {
		log.Error().Err(err).Str("field", field).Msg("Error parsing sun time")
		panic(err)
	}
	return tm
}

// GetSunriseToday returns today's sunrise time.
func GetSunriseToday() time.Time { return getTodayTime("sunrise") }

// GetSunsetToday returns today's sunset time.
func GetSunsetToday() time.Time { return getTodayTime("sunset") }
