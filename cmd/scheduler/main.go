package main

import (
	"context"
	"os"

	"github.com/mikahozz/gohome/integrations/shelly"
	"github.com/mikahozz/gohome/integrations/sun"
	"github.com/rs/zerolog/log"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error().
				Interface("panic", r).
				Stack().
				Msg("Fatal panic occurred")
			os.Exit(1)
		}
	}()

	scheduler := NewScheduler()
	scheduler.AddSchedule(&DailySchedule{
		Name:     "Night lights ON",
		Category: "night_lights",
		Trigger: Trigger{
			Time: sun.GetSunsetToday,
		},
		Action: func(ctx context.Context) error { return shelly.TurnOn(ctx) },
	})

	scheduler.AddSchedule(&DailySchedule{
		Name:     "Night lights OFF",
		Category: "night_lights",
		Trigger: Trigger{
			Time: sun.GetSunriseToday,
		},
		Action: func(ctx context.Context) error { return shelly.TurnOff(ctx) },
	})
	scheduler.Start()
	defer scheduler.Stop()

	// Keep the main function running
	select {}
}
