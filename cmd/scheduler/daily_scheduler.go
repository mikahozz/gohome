package main

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Clock interface for dependency injection in tests
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

// RealClock uses actual system time
type RealClock struct{}

func (rc *RealClock) Now() time.Time {
	return time.Now()
}

func (rc *RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

type TriggerType string

const (
	TriggerTime TriggerType = "time"
)

type FilterType string

const (
	FilterDate FilterType = "date"
)

type AndOrType string

const (
	AND AndOrType = "and"
	OR  AndOrType = "or"
)

type Trigger struct {
	Type TriggerType
	Time func() time.Time
}

type Comparator string

const (
	LessThan    Comparator = "less_than"
	GreaterThan Comparator = "greater_than"
	Equal       Comparator = "equal"
)

type Filter struct {
	Type       FilterType
	Date       time.Time
	Comparator Comparator
}

type Schedule struct {
	Name          string
	Trigger       Trigger
	FilterLogic   AndOrType
	Filters       []Filter
	Action        func(context.Context)
	LastTriggered time.Time
}

type Scheduler struct {
	schedules []*Schedule
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	clock     Clock
}

// NewScheduler creates a new scheduler instance with real clock
func NewScheduler() *Scheduler {
	return NewSchedulerWithClock(&RealClock{})
}

// NewSchedulerWithClock creates a new scheduler instance with custom clock
func NewSchedulerWithClock(clock Clock) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		schedules: make([]*Schedule, 0),
		ctx:       ctx,
		cancel:    cancel,
		clock:     clock,
	}
}

// AddSchedule adds a schedule to the scheduler
func (s *Scheduler) AddSchedule(schedule *Schedule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.schedules = append(s.schedules, schedule)
}

// Start begins running the scheduler
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.clock.After(1 * time.Minute):
			s.evaluate(s.clock.Now())
		}
	}
}

// evaluate checks all schedules and executes matching ones
func (s *Scheduler) evaluate(now time.Time) {
	s.mu.RLock()
	schedules := make([]*Schedule, len(s.schedules))
	copy(schedules, s.schedules)
	s.mu.RUnlock()

	for _, schedule := range schedules {
		if s.shouldTrigger(schedule, now) && s.filtersPass(schedule, now) {
			go schedule.Action(s.ctx)
			schedule.LastTriggered = now
		}
	}
}

// shouldTrigger checks if the trigger condition is met
func (s *Scheduler) shouldTrigger(schedule *Schedule, now time.Time) bool {
	switch schedule.Trigger.Type {
	case TriggerTime:
		t := schedule.Trigger.Time()
		// Check if current time is at or after trigger time
		// If hour is greater, trigger. If hour is equal, check minutes.
		ret := !hasTriggeredThisPeriod(schedule, now) &&
			(now.Hour() > t.Hour() || (now.Hour() == t.Hour() && now.Minute() >= t.Minute()))
		return ret
	}
	return false
}

func hasTriggeredThisPeriod(schedule *Schedule, now time.Time) bool {
	if schedule.LastTriggered.IsZero() {
		return false
	}
	return schedule.LastTriggered.Year() == now.Year() &&
		schedule.LastTriggered.Month() == now.Month() &&
		schedule.LastTriggered.Day() == now.Day()
}

// filtersPass checks if all filters pass according to logic type
func (s *Scheduler) filtersPass(schedule *Schedule, now time.Time) bool {
	if len(schedule.Filters) == 0 {
		return true
	}

	if schedule.FilterLogic == OR {
		for _, filter := range schedule.Filters {
			if s.filterPass(filter, now) {
				return true
			}
		}
		return false
	}

	// Default to AND logic
	for _, filter := range schedule.Filters {
		if !s.filterPass(filter, now) {
			return false
		}
	}
	return true
}

// filterPass checks if a single filter passes
func (s *Scheduler) filterPass(filter Filter, now time.Time) bool {
	switch filter.Type {
	case FilterDate:
		switch filter.Comparator {
		case Equal:
			return now.Year() == filter.Date.Year() &&
				now.Month() == filter.Date.Month() &&
				now.Day() == filter.Date.Day()
		case LessThan:
			return now.Before(filter.Date)
		case GreaterThan:
			return now.After(filter.Date)
		default:
			log.Info().Msg("No filter matched for: " + filter.Date.String())
		}
	}
	return true
}
