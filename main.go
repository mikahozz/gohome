package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mikaahopelto/gohome/integrations/fmi"
	"github.com/mikaahopelto/gohome/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const port = ":9999"

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

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	mux := http.NewServeMux()
	mux.HandleFunc("/weathernow", jsonResponse(fmi.GetWeatherData))
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(mock.IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", jsonResponse(mock.OutdoorWeatherFore))
	mux.HandleFunc("/electricity/prices", jsonResponse(mock.ElectricityPrices))
	log.Fatal().Err(http.ListenAndServe(port, mux))
}
