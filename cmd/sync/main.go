package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/mikahozz/gohome/config"
	"github.com/mikahozz/gohome/dbsync"
	"github.com/mikahozz/gohome/integrations/spot"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}).With().Timestamp().Logger()

	config.LoadEnv()

	db, err := dbsync.SetupDbConn()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to setup database connection")
	}
	defer db.Close()

	store := dbsync.NewMeasurementStore(db)

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
					logger.Info().
						Str("start", start.Format("2006-01-02")).
						Str("end", end.Format("2006-01-02")).
						Msg("No price data available yet")
				} else {
					logger.Error().Err(err).Msg("Error fetching prices")
				}
				logger.Info().Msg("Sleeping for 10 minutes")
				time.Sleep(10 * time.Minute)
				continue
			}

			jsonData, err := json.MarshalIndent(spotPriceList, "", "  ")
			if err != nil {
				logger.Error().Err(err).Msg("Error marshalling prices to JSON")
				continue
			}

			logger.Info().
				Str("start", start.Format("2006-01-02")).
				Str("end", end.Format("2006-01-02")).
				RawJSON("prices", jsonData).
				Msg("Got prices")

			// Store prices in database
			for _, price := range spotPriceList.Prices {
				measurement := &dbsync.Measurement{
					Timestamp: price.DateTime,
					SensorID:  "spot_price",
					MainValue: price.PriceCkwh,
					Value: map[string]interface{}{
						"price": price.PriceCkwh,
						"unit":  "EUR/MWh",
					},
				}

				if err := store.Insert(ctx, measurement); err != nil {
					logger.Error().Err(err).
						Time("timestamp", price.DateTime).
						Float64("price", price.PriceCkwh).
						Msg("Failed to store price")
					continue
				}

				logger.Debug().
					Time("timestamp", price.DateTime).
					Float64("price", price.PriceCkwh).
					Msg("Stored price")
			}

			break
		}

		nextPollTime := time.Date(now.Year(), now.Month(), now.Day()+1, 12, 0, 0, 0, loc)
		logger.Info().Time("next_poll", nextPollTime).Msg("Scheduled next poll")
		time.Sleep(time.Until(nextPollTime))
	}
}
