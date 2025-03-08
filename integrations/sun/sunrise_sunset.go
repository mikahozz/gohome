package sun

import (
	"encoding/json"
	"io/ioutil"
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
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var sunData SunData
	if err := json.Unmarshal(data, &sunData); err != nil {
		return nil, err
	}

	return &sunData, nil
}

// GetDailyData returns the sun data for a specific date
func (s *SunData) GetDailyData(date time.Time) *DailyData {
	dateStr := date.Format("2006-01-02")

	for _, daily := range s.Results {
		if daily.Date == dateStr {
			return &daily
		}
	}

	return nil
}
