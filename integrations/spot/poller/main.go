package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mikahozz/gohome/integrations/spot"
)

func main() {
	for {
		now := time.Now()
		loc, _ := time.LoadLocation("UTC")
		tomorrow := now.AddDate(0, 0, 1)
		start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)

		for {
			spotPriceList, err := spot.GetPrices(start, end, loc)
			if err != nil {
				if _, isNoData := err.(*spot.NoDataError); isNoData {
					fmt.Printf("No price data available yet for period %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
				} else {
					fmt.Printf("Error fetching prices: %v\n", err)
				}
				fmt.Println("Sleeping for 10 minutes")
				time.Sleep(10 * time.Minute)
				continue
			}

			jsonData, err := json.MarshalIndent(spotPriceList, "", "  ")
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				return
			}
			fmt.Printf("Got prices for %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
			fmt.Println(string(jsonData))

			break
		}

		nextPollTime := time.Date(now.Year(), now.Month(), now.Day()+1, 12, 0, 0, 0, loc)
		fmt.Println("Next poll time:", nextPollTime)
		time.Sleep(time.Until(nextPollTime))
	}
}
