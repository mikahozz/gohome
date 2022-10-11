package fmi

import (
	"fmt"
	"net/http"
)

const port = ":9999"

func jsonResponse(data string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, data)
	}
}

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/weathernow", jsonResponse("{testi: 'sdf'}"))
	// log.Fatal(http.ListenAndServe(port, mux))
}

//func GetObservations() {
//fmiEndpoint := "http://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=getFeature&storedquery_id=fmi::observations::weather::hourly::multipointcoverage&place=kumpula"
//observationsRaw, err := http.Get(fmiEndpoint)
// 	if err != nil {
// 		log.Print("Error fetching observations from FMI", err)
// 	}
// }
