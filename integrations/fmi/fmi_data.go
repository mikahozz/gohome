package fmi

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type FeatureCollection struct {
	GridSeriesObservation GridSeriesObservation `xml:"member>GridSeriesObservation"`
	Resolution            Resolution
}
type GridSeriesObservation struct {
	BeginPosition              string  `xml:"phenomenonTime>TimePeriod>beginPosition"`
	EndPosition                string  `xml:"phenomenonTime>TimePeriod>endPosition"`
	DoubleOrNilReasonTupleList string  `xml:"result>MultiPointCoverage>rangeSet>DataBlock>doubleOrNilReasonTupleList"`
	Fields                     []Field `xml:"result>MultiPointCoverage>rangeType>DataRecord>field"`
}
type Field struct {
	Name string `xml:"name,attr"`
}
type Resolution int64

const (
	Hours Resolution = iota + 1
	Minutes
)

func GetData(location string, r Resolution) {

}

func (f FeatureCollection) ConvertToWeatherData() []WeatherData {
	if f.Resolution == 0 {
		log.Fatal("Resolution is not set, cannot convert to WeatherData")
	}
	wArr := []WeatherData{}
	lines := strings.Split(
		strings.TrimSpace(
			strings.ReplaceAll(f.GridSeriesObservation.DoubleOrNilReasonTupleList, "\r\n", "\n"),
		),
		"\n")
	beginDate, err := time.Parse(time.RFC3339, f.GridSeriesObservation.BeginPosition)
	if err != nil {
		log.Fatalf("Failed to parse date: %s", f.GridSeriesObservation.BeginPosition)
	}
	dt := beginDate
	var timeAdd time.Duration
	if f.Resolution == Hours {
		timeAdd = time.Hour
	}
	if f.Resolution == Minutes {
		timeAdd = time.Minute * 10
	}
	for i, line := range lines {
		w := WeatherData{}
		w.Time = dt.UTC().Format(time.RFC3339)
		values := strings.Split(strings.TrimSpace(line), " ")
		fields := f.GridSeriesObservation.Fields
		if len(values) != len(fields) {
			log.Fatalf("measures len: %d != fields len: %d", len(values), len(fields))
		}
		for j, field := range fields {
			value, err := strconv.ParseFloat(values[j], 64)
			if err != nil {
				log.Fatalf("Failed to parse string measure %s from position %d from line %d: %v", values[j], j, i, err)
			}
			switch field.Name {
			case "TA_PT1H_AVG", "t2m":
				w.Temp = valueOrZero(value)
			case "TA_PT1H_MAX":
				w.TempMax = valueOrZero(value)
			case "TA_PT1H_MIN":
				w.TempMin = valueOrZero(value)
			case "RH_PT1H_AVG", "rh":
				w.Humidity = valueOrZero(value)
			case "WS_PT1H_AVG", "ws_10min":
				w.WindSpeed = valueOrZero(value)
			case "WS_PT1H_MAX", "wg_10min":
				w.MaxWindSpeed = valueOrZero(value)
			case "WS_PT1H_MIN":
				w.MinWindSpeed = valueOrZero(value)
			case "WD_PT1H_AVG", "wd_10min":
				w.WindDirection = valueOrZero(value)
			case "PRA_PT1H_ACC", "r_1h":
				w.Rain = valueOrZero(value)
			case "PRI_PT1H_MAX", "ri_10min":
				w.MaxRainIntensity = valueOrZero(value)
			case "PA_PT1H_AVG", "p_sea":
				w.Pressure = valueOrZero(value)
			case "WAWA_PT1H_RANK", "wawa":
				w.Weather = valueOrZero(value)
			case "td":
				w.DewPoint = valueOrZero(value)
			case "snow_aws":
				w.SnowDepth = valueOrZero(value)
			case "vis":
				w.Visibility = valueOrZero(value)
			case "n_man":
				w.CloudCover = valueOrZero(value)
			}
		}
		wArr = append(wArr, w)
		dt = dt.Add(timeAdd)
	}
	return wArr
}

func valueOrZero(v float64) float64 {
	if math.IsNaN(v) {
		return 0.0
	}
	return v
}
