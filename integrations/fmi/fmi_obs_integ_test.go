package fmi

import (
	"testing"
)

func TestGetWeatherData(t *testing.T) {
	obs := FMI_ObservationsModel{}
	err := obs.Get_FMIObservations(StationId("101004"))
	if err != nil {
		t.Errorf("Get_FMIObservations failed: %v", err)
	}
	_, err = obs.ConvertToWeatherData()
	if err != nil {
		t.Errorf("ConvertToWeatherData failed: %v", err)
	}

}
