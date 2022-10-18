package mock

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":9999"

func jsonResponse(f func() (string, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json, err := f()
		if err != nil {
			http.Error(w, "Error occurred in mock endpoint ", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, json)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/indoor/dev_upstairs", jsonResponse(IndoorDevUpstairs))
	mux.HandleFunc("/weatherfore", jsonResponse(OutdoorWeatherFore))
	mux.HandleFunc("/weathernow", jsonResponse(OutdoorWeathernNow))
	mux.HandleFunc("/electricity/prices", jsonResponse(ElectricityPrices))
	log.Fatal(http.ListenAndServe(port, mux))
}
