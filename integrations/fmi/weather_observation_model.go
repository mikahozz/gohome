package main

type FeatureCollection struct {
	GridSeriesObservation GridSeriesObservation `xml:"member>GridSeriesObservation"`
}
type GridSeriesObservation struct {
	BeginPosition string `xml:"phenomenonTime>TimePeriod>beginPosition"`
	EndPosition   string `xml:"phenomenonTime>TimePeriod>endPosition"`
}
