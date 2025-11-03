package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createSchedule(t func() time.Time, fn func(context.Context)) *Scheduler {
	scheduler := NewScheduler()

	schedule := &Schedule{
		Name: "Test Schedule",
		Trigger: Trigger{
			Type: TriggerTime,
			Time: t,
		},
		Action: fn,
	}
	scheduler.AddSchedule(schedule)

	return scheduler
}

func TestSunsetSchedule(t *testing.T) {
	// Track if action was called
	var actionCalled bool
	var mu sync.Mutex

	dummyAction := func(ctx context.Context) {
		mu.Lock()
		actionCalled = true
		mu.Unlock()
	}

	t.Run("Should execute at 18:30", func(t *testing.T) {
		now := time.Now()
		actionCalled = false
		scheduleTime := func() time.Time {
			return time.Date(now.Year(), now.Month(), now.Day(), 18, 30, 0, 0, time.Local)
		}
		testTime := time.Date(now.Year(), now.Month(), now.Day(), 18, 30, 0, 0, time.Local)
		scheduler := createSchedule(scheduleTime, dummyAction)
		scheduler.evaluate(testTime)

		time.Sleep(100 * time.Millisecond)

		assert.True(t, actionCalled, "Action should have been executed at 18:30")
	})

	t.Run("Should execute at 18:50", func(t *testing.T) {
		now := time.Now()
		actionCalled = false
		scheduleTime := func() time.Time {
			return time.Date(now.Year(), now.Month(), now.Day(), 18, 30, 0, 0, time.Local)
		}
		testTime := time.Date(now.Year(), now.Month(), now.Day(), 18, 50, 0, 0, time.Local)
		scheduler := createSchedule(scheduleTime, dummyAction)
		scheduler.evaluate(testTime)

		time.Sleep(100 * time.Millisecond)

		assert.True(t, actionCalled, "Action should have been executed at 18:50")
	})

	t.Run("Should not execute at 18:00", func(t *testing.T) {
		now := time.Now()
		actionCalled = false
		scheduleTime := func() time.Time {
			return time.Date(now.Year(), now.Month(), now.Day(), 18, 30, 0, 0, time.Local)
		}
		testTime := time.Date(now.Year(), now.Month(), now.Day(), 18, 00, 0, 0, time.Local)
		scheduler := createSchedule(scheduleTime, dummyAction)
		scheduler.evaluate(testTime)

		time.Sleep(100 * time.Millisecond)

		assert.False(t, actionCalled, "Action should not have been executed at 18:00")
	})
	t.Run("Should not execute second time same day but day after", func(t *testing.T) {
		now := time.Now()
		actionCalled = false
		scheduleTime := func() time.Time {
			return time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, time.Local)
		}
		testTime := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, time.Local)
		scheduler := createSchedule(scheduleTime, dummyAction)
		scheduler.evaluate(testTime)

		time.Sleep(100 * time.Millisecond)

		assert.True(t, actionCalled, "Action should have been executed first time")

		actionCalled = false
		testTime = time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, time.Local)
		scheduler.evaluate(testTime)

		time.Sleep(100 * time.Millisecond)

		assert.False(t, actionCalled, "Action should NOT have been executed second time same day")

		tomorrowTestTime := testTime.Add(24 * time.Hour)
		scheduler.evaluate(tomorrowTestTime)

		time.Sleep(100 * time.Millisecond)

		assert.True(t, actionCalled, "Action should have been executed next day")
	})
}
