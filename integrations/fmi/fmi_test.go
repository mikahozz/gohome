package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"strings"
	"testing"
)

func TestXml(t *testing.T) {
	var featureCollection = &FeatureCollection{}

	xmlObservations, err := os.Open("example.xml")
	if err != nil {
		t.Fatal("Could not retrieve example.xml file", err)
	}
	defer xmlObservations.Close()
	decoder := xml.NewDecoder(xmlObservations)
	err = decoder.Decode(featureCollection)
	if err != nil {
		t.Fatal("Could not decode example.xml file", err)
	}
	//log.Printf("%+v", featureCollection)
	t.Run("Assert dates", func(t *testing.T) {
		AssertEqual(featureCollection.GridSeriesObservation.BeginPosition, "2022-10-01T07:00:00Z", t)
		AssertEqual(featureCollection.GridSeriesObservation.EndPosition, "2022-10-02T07:00:00Z", t)
	})
	if len(featureCollection.GridSeriesObservation.Fields) != 12 {
		t.Errorf("GridSeriesObservation.Fields != 12")
	}
	if len(strings.TrimSpace(featureCollection.GridSeriesObservation.DoubleOrNilReasonTupleList)) < 56*25 {
		t.Errorf("DoubleOrNilReasonTupleList < 56*25")
	}
	observations := CreateObservationsModel(*featureCollection)
	json, err := json.Marshal(observations)
	if err != nil {
		t.Errorf("Failed to marshal json from: %+v", observations)
	}
	log.Print(string(json))
	if len(observations) < 25 {
		t.Errorf("observations len < 25, got %d", len(observations))
	}
	AssertEqual(observations[0].DateTime, "2022-10-01T07:00:00Z", t)

}

func AssertEqual(got string, want string, t *testing.T) {
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
