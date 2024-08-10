package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mikaahopelto/gohome/integrations/fmi"
	"github.com/mikaahopelto/gohome/mock"
	"github.com/mikahozz/gohome/integrations/cal"
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

func getWeatherData(place string, requestType fmi.RequestType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weather, err := fmi.GetWeatherData(fmi.StationId(place), requestType)
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
func getCalendarEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := cal.DateOffset{}
		to := cal.DateOffset{Days: 7}
		events, err := cal.GetFamilyCalendarEvents(from, to)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, fmt.Sprintf("Error occurred fetching calendar events"), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(events)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, fmt.Sprintf("Error occurred in json conversion of calendar events"), http.StatusInternalServerError)
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
	mux.HandleFunc("/weathernow", getWeatherData("101004", fmi.Observations))
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(mock.IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", getWeatherData("Tapanila,Helsinki", fmi.Forecast))
	mux.HandleFunc("/electricity/prices", jsonResponse(mock.ElectricityPrices))
	mux.HandleFunc("/api/events", getCalendarEvents())
	log.Fatal().Err(http.ListenAndServe(port, mux))
}
