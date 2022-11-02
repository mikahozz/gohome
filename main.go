package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mikaahopelto/gohome/integrations/fmi"
	"github.com/mikaahopelto/gohome/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const port = ":6001"

func jsonResponse(f func() (string, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json, err := f()
		if err != nil {
			log.Error().Err(err).Msg("")
			http.Error(w, "Error occurred when performing request", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, json)
	}
}

func getWeatherData(place string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weather, err := fmi.GetWeatherData(fmi.StationId(place))
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, fmt.Sprintf("Error occurred in fetching weather data for %s", place), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(weather.WeatherData)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, fmt.Sprintf("Error occurred in fetching weather data for %s", place), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	mux := http.NewServeMux()
	mux.HandleFunc("/weathernow", getWeatherData("101004"))
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(mock.IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", jsonResponse(mock.OutdoorWeatherFore))
	mux.HandleFunc("/electricity/prices", jsonResponse(mock.ElectricityPrices))
	log.Fatal().Err(http.ListenAndServe(port, mux))
}
