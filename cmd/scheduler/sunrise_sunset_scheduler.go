package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mikahozz/gohome/integrations/sun"
	"github.com/rs/zerolog/log"
)

func nightLightsOn(ctx context.Context) {
	log.Info().Msg("Evening lights turned ON")
}
func nightLightsOff(ctx context.Context) {
	log.Info().Msg("Evening lights turned OFF")
}

func CreateSunriseSunsetScheduler() *Scheduler {
	// Load sun data - note that only month and day are considered, year is ignored
	sunData, err := sun.LoadSunData("../../integrations/sun/sun_helsinki_2025.json")
	if err != nil {
		log.Error().Err(err).Msg("Error loading sun data")
		panic(err)
	}

	sunriseSunset := func(t time.Time) (time.Time, time.Time) {
		sArr := sunData.GetDailyData(time.Now(), nil)
		if len(sArr) == 0 {
			log.Error().Err(err).Msg("No sun data found for the date")
			panic(err)
		}
		s := sArr[0]
		sunrise, err := time.Parse("2006-01-02 3:04:05 PM", fmt.Sprintf("%s %s", s.Date, s.Sunrise))
		if err != nil {
			log.Error().Err(err).Msg("Error parsing sunrise time")
			panic(err)
		}
		sunset, err := time.Parse("2006-01-02 3:04:05 PM", fmt.Sprintf("%s %s", s.Date, s.Sunset))
		if err != nil {
			log.Error().Err(err).Msg("Error parsing sunset time")
			panic(err)
		}
		return sunrise, sunset
	}
	sunrise := func() time.Time {
		sunrise, _ := sunriseSunset(time.Now())
		return sunrise
	}
	sunset := func() time.Time {
		_, sunset := sunriseSunset(time.Now())
		return sunset
	}

	scheduler := NewScheduler()
	eveningSchedule := &Schedule{
		Name: "Night Lights ON",
		Trigger: Trigger{
			Type: TriggerTime,
			Time: sunset,
		},
		Action: nightLightsOn,
	}
	scheduler.AddSchedule(eveningSchedule)

	morningSchedule := &Schedule{
		Name: "Night Lights OFF",
		Trigger: Trigger{
			Type: TriggerTime,
			Time: sunrise,
		},
		Action: nightLightsOff,
	}
	scheduler.AddSchedule(morningSchedule)

	return scheduler
}
