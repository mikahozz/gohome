package main

import (
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
		Name: "Night lights ON",
		Trigger: Trigger{
			Time: sun.GetSunsetToday,
		},
		Action: shelly.TurnOn,
	})

	scheduler.AddSchedule(&DailySchedule{
		Name: "Night lights OFF",
		Trigger: Trigger{
			Time: sun.GetSunriseToday,
		},
		Action: shelly.TurnOff,
	})
	scheduler.Start()
	defer scheduler.Stop()

	// Keep the main function running
	select {}
}
