package fmi

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func LoadXml(t *testing.T, fn string, fc *FeatureCollection, r Resolution) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Fatalf("Could not retrieve %s file: %v", fn, err)
	}
	err = xml.Unmarshal(data, fc)
	if err != nil {
		t.Fatalf("Could not retrieve %s file: %v", fn, err)
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
	LastObservationTime  string
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
		LastObservationTime:  "2022-10-10T14:50:00Z",
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
		LastObservationTime:  "2022-10-02T07:00:00Z",
		Resolution:           Hours,
	}
	WeatherDataTests(t, v)
}

func TestInvalidXml(t *testing.T) {
	var fc = &FeatureCollection{}
	LoadXml(t, "exampleEmpty.xml", fc, Minutes)
	log.Printf("%+v", fc)
	LoadXml(t, "exampleInvalid.xml", fc, Minutes)
	log.Printf("%+v", fc)
}

func WeatherDataTests(t *testing.T, test TestValues) {
	var fc = &FeatureCollection{}
	// Initialize xml
	if test.Resolution == Minutes {
		LoadXml(t, "exampleMinutes.xml", fc, Minutes)
	} else if test.Resolution == Hours {
		LoadXml(t, "exampleHours.xml", fc, Hours)
	} else {
		t.Error("Missing testValues.Resolution")
	}
	//log.Printf("%+v", featureCollection)

	if fc.GridSeriesObservation.BeginPosition != test.BeginPosition {
		t.Errorf("BeginPosition, got %s, want %s", fc.GridSeriesObservation.BeginPosition, test.BeginPosition)
	}
	if fc.GridSeriesObservation.EndPosition != test.EndPosition {
		t.Errorf("EndPosition, got %s, want %s", fc.GridSeriesObservation.EndPosition, test.EndPosition)
	}
	if len(fc.GridSeriesObservation.Fields) != test.FieldsLen {
		t.Errorf("GridSeriesObservation.Fields len, got %d, want %d", len(fc.GridSeriesObservation.Fields), test.FieldsLen)
	}
	if len(strings.TrimSpace(fc.GridSeriesObservation.DoubleOrNilReasonTupleList)) < test.TupleListMinLen {
		t.Errorf("DoubleOrNilReasonTupleList min len, got %d, want %d", len(strings.TrimSpace(fc.GridSeriesObservation.DoubleOrNilReasonTupleList)), test.TupleListMinLen)
	}
	// Load xml into WeatherData
	weather := fc.ConvertToWeatherData()
	_, err := json.Marshal(weather)
	if err != nil {
		t.Errorf("Failed to marshal json from: %+v. Err: %v", weather, err)
	}
	// log.Print(string(json))
	if len(weather) != test.ObservationsLen {
		t.Errorf("len(observations), got %d, want %d", len(weather), test.ObservationsLen)
	}
	if weather[0].Time != test.FirstObservationTime {
		t.Errorf("observations[0].Time, got %s, want %s", weather[0].Time, test.FirstObservationTime)
	}
	if weather[len(weather)-1].Time != test.LastObservationTime {
		t.Errorf("last weather time != LastObservationTime, got %s, want %s", weather[len(weather)-1].Time, test.LastObservationTime)
	}
}
