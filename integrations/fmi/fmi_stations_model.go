package fmi

type FMI_StationsModel struct {
	Stations StationCollection `xml:"FeatureCollection"`
}

type StationCollection struct {
	Stations []Station `xml:"member>EnvironmentalMonitoringFacility"` // Weather stations
}
type Station struct {
	Id    string `xml:"identifier"`
	Names []Name `xml:"name"`
	Point string `xml:"representativePoint>Point>pos"`
}
type Name struct {
	Key   string `xml:"codeSpace,attr"`
	Value string `xml:",chardata"`
}
