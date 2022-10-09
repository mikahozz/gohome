package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type FeatureCollection struct {
	GridSeriesObservation GridSeriesObservation `xml:"member>GridSeriesObservation"`
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

func CreateObservationsModel(f FeatureCollection) []Observation {
	observations := []Observation{}
	measureLines := strings.Split(
		strings.TrimSpace(
			strings.ReplaceAll(f.GridSeriesObservation.DoubleOrNilReasonTupleList, "\r\n", "\n"),
		),
		"\n")
	beginDate, err := time.Parse(time.RFC3339, f.GridSeriesObservation.BeginPosition)
	if err != nil {
		log.Fatalf("Failed to parse date: %s", f.GridSeriesObservation.BeginPosition)
	}
	for i, measureLine := range measureLines {
		observation := Observation{}
		observation.DateTime = beginDate.Add(time.Hour * time.Duration(i)).UTC().Format(time.RFC3339)
		measures := strings.Split(strings.TrimSpace(measureLine), " ")
		fields := f.GridSeriesObservation.Fields
		if len(measures) != len(fields) {
			log.Fatalf("measures len: %d != fields len: %d", len(measures), len(fields))
		}
		for j, field := range fields {
			value, err := strconv.ParseFloat(measures[j], 64)
			if err != nil {
				log.Fatalf("Failed to parse string measure %s from position %d from line %d: %v", measures[j], j, i, err)
			}
			switch field.Name {
			case "TA_PT1H_AVG":
				observation.AirTemperature = value
			case "TA_PT1H_MAX":
				observation.HighestTemperature = value
			case "TA_PT1H_MIN":
				observation.LowestTemperature = value
			case "RH_PT1H_AVG":
				observation.RelativeHumidity = value
			case "WS_PT1H_AVG":
				observation.WindSpeed = value
			case "WS_PT1H_MAX":
				observation.MaximumWindSpeed = value
			case "WS_PT1H_MIN":
				observation.MinimumWindSpeed = value
			case "WD_PT1H_AVG":
				observation.WindDirection = value
			case "PRA_PT1H_ACC":
				observation.PrecipitationAmount = value
			case "PRI_PT1H_MAX":
				observation.MaximumPrecipitationIntensity = value
			case "PA_PT1H_AVG":
				observation.AirPressure = value
			case "WAWA_PT1H_RANK":
				observation.PresentWeather = value
			}
		}
		observations = append(observations, observation)
	}
	return observations
}

// fieldDescriptions := map[FieldTypes]string{
// 	TA_PT1H_AVG:    "Air Temperature",
// 	TA_PT1H_MAX:    "Highest temperature",
// 	TA_PT1H_MIN:    "Lowest temperature",
// 	RH_PT1H_AVG:    "Relative humidity",
// 	WS_PT1H_AVG:    "Wind speed",
// 	WS_PT1H_MAX:    "Maximum wind speed",
// 	WS_PT1H_MIN:    "Minimum wind speed",
// 	WD_PT1H_AVG:    "Wind direction",
// 	PRA_PT1H_ACC:   "Precipitation amount",
// 	PRI_PT1H_MAX:   "Maximum precipitation intensity",
// 	PA_PT1H_AVG:    "Air pressure",
// 	WAWA_PT1H_RANK: "Present weather",
// }

// type FieldTypes int

// const (
// 	TA_PT1H_AVG FieldTypes = iota
// 	TA_PT1H_MAX
// 	TA_PT1H_MIN
// 	RH_PT1H_AVG
// 	WS_PT1H_AVG
// 	WS_PT1H_MAX
// 	WS_PT1H_MIN
// 	WD_PT1H_AVG
// 	PRA_PT1H_ACC
// 	PRI_PT1H_MAX
// 	PA_PT1H_AVG
// 	WAWA_PT1H_RANK
// )
