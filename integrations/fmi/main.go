package fmi

// Param: place: Name of a weather station in https://en.ilmatieteenlaitos.fi/observation-stations?filterKey=groups&filterQuery=weather
// func GetWeatherData(place string) (string, error) {
// 	stationsXml := GetStations()
// 	q := fmt.Sprintf("http://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=getFeature&storedquery_id=fmi::observations::weather::hourly::multipointcoverage&place=%s",
// 		place)
// 	_, err := http.Get(q)
// 	if err != nil {
// 		return "", errors.Wrap(err, "Error fetching data from FMI")
// 	}
// 	return "lkj", nil
// }

// func GetStations() ([]WeatherStation, error) {
// 	url := "https://opendata.fmi.fi/wfs/fin?request=getFeature&storedquery_id=fmi::ef::stations"
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "Could not get: %s: %v", url)
// 	}
// 	defer response.Body.Close()
// 	data, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Could not read response.Body")
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		return nil, errors.Errorf("Error in fetching weather stations. Statuscode: %d. Body: %v", response.StatusCode, data)
// 	}
// 	x := response.Body.Read()
// 	f := &FeatureCollection{}
// 	err = xml.Unmarshal(data.Body.Read(), f)
// 	if err != nil {
// 		return nil, errors.Wrap("Could not parse xml: %v", err)
// 	}
// 	return ConvertToWeatherStations(f)
// }
