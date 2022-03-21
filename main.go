package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mikaahopelto/gohome/mock"
)

const port = ":9999"

func jsonResponse(mockFunc func() string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, mockFunc())
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(mock.IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", jsonResponse(mock.OutdoorWeatherFore))
	mux.HandleFunc("/weathernow", jsonResponse(mock.OutdoorWeathernNow))
	mux.HandleFunc("/electricity/prices", jsonResponse(mock.ElectricityPrices))
	log.Fatal(http.ListenAndServe(port, mux))
}
