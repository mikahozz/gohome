package spot

import (
	"time"
)

// SpotPrice represents a single price point at a specific time
type SpotPrice struct {
	DateTime  time.Time
	PriceCkwh float64
}

// SpotPriceList is a collection of SpotPrice entries
type SpotPriceList struct {
	Prices []SpotPrice
}
