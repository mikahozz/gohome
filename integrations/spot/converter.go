package spot

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseISO8601Duration(duration string) (time.Duration, error) {
	// Remove the "PT" prefix
	duration = strings.TrimPrefix(duration, "PT")

	// Check if the duration ends with "M" for minutes
	if strings.HasSuffix(duration, "M") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(duration, "M"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration format: %s", duration)
		}
		return time.Duration(minutes) * time.Minute, nil
	}

	// Check if the duration ends with "H" for hours
	if strings.HasSuffix(duration, "H") {
		hours, err := strconv.Atoi(strings.TrimSuffix(duration, "H"))
		if err != nil {
			return 0, fmt.Errorf("invalid duration format: %s", duration)
		}
		return time.Duration(hours) * time.Hour, nil
	}

	return 0, fmt.Errorf("unsupported duration format: %s", duration)
}

func ConvertToSpotPriceList(doc *PublicationMarketDocument, periodStart, periodEnd time.Time) (*SpotPriceList, error) {
	var spotPrices []SpotPrice

	for _, ts := range doc.TimeSeries {
		start, err := time.Parse("2006-01-02T15:04Z", ts.Period.TimeInterval.Start)
		if err != nil {
			return nil, fmt.Errorf("error parsing start time: %w", err)
		}
		start = start.UTC() // Ensure we're working with UTC

		resolution, err := parseISO8601Duration(ts.Period.Resolution)
		if err != nil {
			return nil, fmt.Errorf("error parsing resolution: %w", err)
		}

		for _, point := range ts.Period.Points {
			dateTime := start.Add(time.Duration(point.Position-1) * resolution)

			if dateTime.Before(periodStart) || dateTime.After(periodEnd) {
				continue
			}

			// Convert price from EUR/MWh to cents/kWh
			price := point.Price * 100 / 1000

			spotPrices = append(spotPrices, SpotPrice{
				DateTime:  dateTime,
				PriceCkwh: price,
			})
		}
	}

	// Sort the prices by datetime
	sort.Slice(spotPrices, func(i, j int) bool {
		return spotPrices[i].DateTime.Before(spotPrices[j].DateTime)
	})

	fmt.Printf("Total spot prices: %d\n", len(spotPrices))
	if len(spotPrices) > 0 {
		fmt.Printf("First price: %+v\n", spotPrices[0])
		fmt.Printf("Last price: %+v\n", spotPrices[len(spotPrices)-1])
	}

	return &SpotPriceList{Prices: spotPrices}, nil
}
