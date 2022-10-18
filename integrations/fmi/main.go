package fmi

import (
	"net/http"

	"github.com/pkg/errors"
)

func GetWeatherData() (string, error) {
	//q := "http://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=getFeature&storedquery_id=fmi::observations::weather::hourly::multipointcoverage&place=kumpula"
	q := "ksdjf"
	_, err := http.Get(q)
	if err != nil {
		return "", errors.Wrap(err, "Error fetching data from FMI")
	}
	return "lkj", nil
}
