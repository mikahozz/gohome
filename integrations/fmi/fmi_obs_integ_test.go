//go:build integration

package fmi

import (
	"testing"
)

func TestGetWeatherData(t *testing.T) {
	obs := FMI_ObservationsModel{}
	err := obs.LoadObservations(StationId("101004"))
	if err != nil {
		t.Fatalf("LoadObservations failed: %v", err)
	}
	_, err = obs.ConvertToWeatherData()
	if err != nil {
		t.Errorf("ConvertToWeatherData failed: %v", err)
	}

}
