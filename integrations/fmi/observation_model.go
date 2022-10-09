package main

type Observation struct {
	DateTime                      string
	AirTemperature                float64
	HighestTemperature            float64
	LowestTemperature             float64
	RelativeHumidity              float64
	WindSpeed                     float64
	MaximumWindSpeed              float64
	MinimumWindSpeed              float64
	WindDirection                 float64
	PrecipitationAmount           float64
	MaximumPrecipitationIntensity float64
	AirPressure                   float64
	PresentWeather                float64
}
