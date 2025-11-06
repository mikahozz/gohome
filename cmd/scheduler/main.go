package main

import (
	"context"
	"os"
	"time"

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
		Name:     "Night lights ON at sunset",
		Category: "night_lights",
		Trigger: Trigger{
			Time: sun.GetSunsetToday,
		},
		Action: func(ctx context.Context) error { return shelly.TurnOn(ctx) },
	})
	scheduler.AddSchedule(&DailySchedule{
		Name:     "Night lights OFF at 23:00",
		Category: "night_lights",
		Trigger: Trigger{
			Time: func() time.Time {
				now := time.Now()
				return time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
			},
		},
		Action: func(ctx context.Context) error { return shelly.TurnOff(ctx) },
	})

	scheduler.AddSchedule(&DailySchedule{
		Name:     "Morning lights ON at 6:45",
		Category: "night_lights",
		Trigger: Trigger{
			Time: func() time.Time {
				now := time.Now()
				return time.Date(now.Year(), now.Month(), now.Day(), 6, 45, 0, 0, now.Location())
			},
		},
		Action: func(ctx context.Context) error { return shelly.TurnOn(ctx) },
	})
	scheduler.AddSchedule(&DailySchedule{
		Name:     "Morning lights OFF at sunrise",
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
