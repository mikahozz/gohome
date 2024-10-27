package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mikahozz/gohome/integrations/cal"
	"github.com/mikahozz/gohome/integrations/fmi"
	"github.com/mikahozz/gohome/integrations/spot"
	"github.com/mikahozz/gohome/mock"
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
		fmt.Fprint(w, json)
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
			http.Error(w, "Error occurred fetching calendar events", http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(events)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, "Error occurred in json conversion of calendar events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func getSpotPrices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")
		timeFormat := r.URL.Query().Get("timeFormat")

		// Default to UTC if timeFormat is not specified
		if timeFormat == "" {
			timeFormat = "utc"
		}

		// Get the location based on timeFormat
		var location *time.Location
		var err error
		switch timeFormat {
		case "utc":
			location = time.UTC
		case "local":
			location = time.Local
		default:
			location, err = time.LoadLocation(timeFormat)
			if err != nil {
				log.Error().Err(err).Msg("Invalid timezone format")
				http.Error(w, "Invalid timezone format", http.StatusBadRequest)
				return
			}
		}

		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, "Invalid start time format. Use RFC3339.", http.StatusBadRequest)
			return
		}

		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			log.Err(err).Msg("")
			http.Error(w, "Invalid end time format. Use RFC3339.", http.StatusBadRequest)
			return
		}

		log.Info().Msgf("Getting spot prices for %s to %s in %s format", start, end, timeFormat)
		prices, err := spot.GetPrices(start, end, location)
		if err != nil {
			log.Error().Err(err).Msg("Error getting spot prices")
			http.Error(w, "Error occurred fetching spot prices", http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(prices)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling spot prices to JSON")
			http.Error(w, "Error occurred in JSON conversion of spot prices", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Starting server")
	mux := http.NewServeMux()
	mux.HandleFunc("/weathernow", getWeatherData("101004", fmi.Observations))
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(mock.IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", getWeatherData("Tapanila,Helsinki", fmi.Forecast))
	mux.HandleFunc("/electricity/prices", getSpotPrices())
	mux.HandleFunc("/api/events", getCalendarEvents())
	log.Fatal().Err(http.ListenAndServe(port, mux)).Send()
}
