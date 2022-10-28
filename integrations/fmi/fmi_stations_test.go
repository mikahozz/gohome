package fmi

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestStations(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/exampleStations.xml")
	if err != nil {
		t.Fatal("Error opening exampleStations.xml file")
	}
	fmi := &FMI_StationsModel{}
	s := &fmi.Stations
	err = xml.Unmarshal(data, s)
	if err != nil {
		t.Fatalf("Could not parse exampleStations.xml: %v", err)
	}
	if l := 452; len(s.Stations) != l {
		t.Errorf("len(f.Stations), got %d, want, %d", len(s.Stations), l)
	}
	if id := "100539"; s.Stations[0].Id != id {
		t.Errorf("Stations[0].Identifier, got %s, want %s", s.Stations[0].Id, id)
	}
	if point := "65.673370 24.515260"; s.Stations[0].Point != point {
		t.Errorf("Stations[0].Point, got %s, want %s", s.Stations[0].Point, point)
	}
	if key := "http://xml.fmi.fi/namespace/locationcode/name"; s.Stations[0].Names[0].Key != key {
		t.Errorf("f.Stations[0].Names[0].Key, got %s, want %s", s.Stations[0].Names[0].Key, key)
	}
	if name := "Kemi Ajos"; s.Stations[0].Names[0].Value != name {
		t.Errorf("f.Stations[0].Names[0].Value, got %s, want %s", s.Stations[0].Names[0].Value, name)
	}
	if id := "874863"; s.Stations[len(s.Stations)-1].Id != id {
		t.Errorf("Stations[last].Identifier, got %s, want %s", s.Stations[0].Id, id)
	}
}
