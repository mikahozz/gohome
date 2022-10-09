package mock

import (
	"encoding/json"
	"time"
)

type price struct {
	DateTime string
	Price    float64
}

func ElectricityPrices() string {
	dates := GenerateFutureDates(time.Hour, 10, false, false)
	var prices []price
	for _, date := range dates {
		prices = append(prices, price{
			DateTime: date,
			Price:    rand.float64() * 20.0,
		})
	}
	pricesJson, err := json.MarshalIndent(prices, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(pricesJson)
}
