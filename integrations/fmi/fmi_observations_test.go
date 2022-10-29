package fmi

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"
)

func LoadXml(t *testing.T, fn string, fc *ObservationCollection, r Resolution) {
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
	weatherDataTests(t, v)
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
	weatherDataTests(t, v)
}

func TestInvalidXml(t *testing.T) {
	fmi := &FMI_ObservationsModel{}
	fc := &fmi.Observations
	LoadXml(t, "testdata/exampleEmpty.xml", fc, Minutes)
	//log.Printf("%+v", fc)
	LoadXml(t, "testdata/exampleInvalid.xml", fc, Minutes)
	//log.Printf("%+v", fc)
}

// func TestAPI(t *testing.T) {
// 	_, err := GetWeatherData("random")
// 	if err == nil {
// 		t.Errorf("Got nil err object, wanted error")
// 	}
// }

func weatherDataTests(t *testing.T, test TestValues) {
	fmi := &FMI_ObservationsModel{}
	fc := &fmi.Observations
	// Initialize xml
	if test.Resolution == Minutes {
		LoadXml(t, "testdata/exampleMinutes.xml", fc, Minutes)
	} else if test.Resolution == Hours {
		LoadXml(t, "testdata/exampleHours.xml", fc, Hours)
	} else {
		t.Error("Missing testValues.Resolution")
	}
	//log.Printf("%+v", featureCollection)
	err := fmi.Validate()
	if err != nil {
		t.Errorf("Error validating observations model: %v", err)
	}
	if fc.Observation.BeginPosition != test.BeginPosition {
		t.Errorf("BeginPosition, got %s, want %s", fc.Observation.BeginPosition, test.BeginPosition)
	}
	if fc.Observation.EndPosition != test.EndPosition {
		t.Errorf("EndPosition, got %s, want %s", fc.Observation.EndPosition, test.EndPosition)
	}
	if len(fc.Observation.Fields) != test.FieldsLen {
		t.Errorf("Observation.Fields len, got %d, want %d", len(fc.Observation.Fields), test.FieldsLen)
	}
	if len(strings.TrimSpace(fc.Observation.Measures)) < test.TupleListMinLen {
		t.Errorf("Measures min len, got %d, want %d", len(strings.TrimSpace(fc.Observation.Measures)), test.TupleListMinLen)
	}
	// Load xml into WeatherData
	weather := fc.ConvertToWeatherData()
	_, err = json.Marshal(weather)
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
