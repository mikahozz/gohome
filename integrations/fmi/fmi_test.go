package fmi

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func LoadXml(t *testing.T, fn string, fc *FeatureCollection, r Resolution) {
	f, err := os.Open(fn)
	if err != nil {
		t.Fatal("Could not retrieve example.xml file", err)
	}
	defer f.Close()
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(fc)
	if err != nil {
		t.Fatal("Could not decode example.xml file", err)
	}
	fc.Resolution = r
}

type TestValues struct {
	BeginPosition        string
	EndPosition          string
	FieldsLen            int
	TupleListMinLen      int
	ObservationsLen      int
	FirstObservationTime string
	Resolution           Resolution
}

func TestWeatherDataMinutes(t *testing.T) {
	v := TestValues{
		BeginPosition:        "2022-10-10T02:50:00Z",
		EndPosition:          "2022-10-10T14:50:00Z",
		FieldsLen:            13,
		TupleListMinLen:      13 * 73 * 4,
		ObservationsLen:      73,
		FirstObservationTime: "2022-10-10T02:50:00Z",
		Resolution:           Minutes,
	}
	WeatherDataTests(t, v)
}

func TestWeatherDataHours(t *testing.T) {
	v := TestValues{
		BeginPosition:        "2022-10-01T07:00:00Z",
		EndPosition:          "2022-10-02T07:00:00Z",
		FieldsLen:            12,
		TupleListMinLen:      12 * 25 * 4,
		ObservationsLen:      25,
		FirstObservationTime: "2022-10-01T07:00:00Z",
		Resolution:           Hours,
	}
	WeatherDataTests(t, v)
}

func WeatherDataTests(t *testing.T, v TestValues) {
	var fc = &FeatureCollection{}
	// Initialize xml
	if v.Resolution == Minutes {
		LoadXml(t, "exampleMinutes.xml", fc, Minutes)
	} else if v.Resolution == Hours {
		LoadXml(t, "exampleHours.xml", fc, Hours)
	} else {
		t.Error("Missing testValues.Resolution")
	}
	//log.Printf("%+v", featureCollection)

	if fc.GridSeriesObservation.BeginPosition != v.BeginPosition {
		t.Errorf("BeginPosition, got %s, want %s", fc.GridSeriesObservation.BeginPosition, v.BeginPosition)
	}
	if fc.GridSeriesObservation.EndPosition != v.EndPosition {
		t.Errorf("EndPosition, got %s, want %s", fc.GridSeriesObservation.EndPosition, v.EndPosition)
	}
	if len(fc.GridSeriesObservation.Fields) != v.FieldsLen {
		t.Errorf("GridSeriesObservation.Fields len, got %d, want %d", len(fc.GridSeriesObservation.Fields), v.FieldsLen)
	}
	if len(strings.TrimSpace(fc.GridSeriesObservation.DoubleOrNilReasonTupleList)) < v.TupleListMinLen {
		t.Errorf("DoubleOrNilReasonTupleList min len, got %d, want %d", len(strings.TrimSpace(fc.GridSeriesObservation.DoubleOrNilReasonTupleList)), v.TupleListMinLen)
	}
	// Load xml into WeatherData
	w := fc.ConvertToWeatherData()
	_, err := json.Marshal(w)
	if err != nil {
		t.Errorf("Failed to marshal json from: %+v. Err: %v", w, err)
	}
	// log.Print(string(json))
	if len(w) != v.ObservationsLen {
		t.Errorf("len(observations), got %d, want %d", len(w), v.ObservationsLen)
	}
	if w[0].Time != v.FirstObservationTime {
		t.Errorf("observations[0].Time, got %s, want %s", w[0].Time, v.FirstObservationTime)
	}
}
