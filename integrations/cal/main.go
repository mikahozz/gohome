package cal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	webdav "github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
)

type Event struct {
	Uid     string    `json:"uid"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Summary string    `json:"summary"`
}

type DateOffset struct {
	Years  int
	Months int
	Days   int
}

// Pretty print DateOffset for logging
func (d DateOffset) String() string {
	return fmt.Sprintf("%d years, %d months, %d days", d.Years, d.Months, d.Days)
}

// GetFamilyCalendarEvents retrieves events from a family calendar within a specified date range.
// The range is determined by two DateOffset structs, 'from' and 'to'.
// Each DateOffset represents an offset in years, months, and days from the current date.
// For example, GetFamilyCalendarEvents(DateOffset{Days: -7}, DateOffset{Days: 7}) retrieves events from one week before to one week after today.
// The function returns a slice of Event structs and an error. If the function succeeds, the error is nil.
// If the function fails, the slice is nil and the error contains details about the failure.
func GetFamilyCalendarEvents(from DateOffset, to DateOffset) ([]Event, error) {
	reqStart := time.Now().AddDate(from.Years, from.Months, from.Days)
	reqEnd := time.Now().AddDate(to.Years, to.Months, to.Days)

	println("Getting family calendar events from: ", from.String(), " to: ", to.String())

	httpClient := &http.Client{}

	fmt.Printf("Connecting to %s with %s\n", config.calUrl, config.username)
	authorizedClient := webdav.HTTPClientWithBasicAuth(httpClient, config.username, config.password)
	calDavClient, err := caldav.NewClient(authorizedClient, config.calUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)

	fmt.Print("Finding current user principal. ")
	curUser, err := calDavClient.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in FindCurrentUserPrincipal: %w", err)
	}
	fmt.Println("Current principal: ", curUser)

	fmt.Print("Finding calendar home. ")
	homeSet, err := calDavClient.FindCalendarHomeSet(ctx, curUser)
	if err != nil {
		return nil, fmt.Errorf("error in FindCalendarHomeSet: %w", err)
	}
	fmt.Println("Calendar home: ", homeSet)

	fmt.Println("Finding calendars. ")
	calendars, err := calDavClient.FindCalendars(ctx, homeSet)
	if err != nil {
		return nil, fmt.Errorf("error in FindCalendars: %w", err)
	}
	for i, cal := range calendars {
		fmt.Printf("Calendar %d: %s: %s\n", i, cal.Path, cal.Name)
	}

	events := []Event{}
	fmt.Printf("Querying calendars with the name '%s'\n", config.calName)
	calQuery := buildQuery(from, to)
	for _, cal := range calendars {
		if cal.Name == config.calName {
			fmt.Printf("Found. Querying calendar: %s\n", cal.Path)
			objects, err := calDavClient.QueryCalendar(ctx, cal.Path, &calQuery)
			if err != nil {
				return nil, fmt.Errorf("error in QueryCalendar: %w", err)
			}
			fmt.Printf("Found %d objects\n", len(objects))
			for _, obj := range objects {
				uid := obj.Data.Children[0].Props.Get("UID")
				dtstart := obj.Data.Children[0].Props.Get("DTSTART")
				dtend := obj.Data.Children[0].Props.Get("DTEND")
				summary := obj.Data.Children[0].Props.Get("SUMMARY")
				var exdateVal string
				exdate := obj.Data.Children[0].Props.Get("EXDATE")
				if exdate != nil {
					exdateVal = exdate.Value
				}
				var rRuleVal string
				rrule := obj.Data.Children[0].Props.Get("RRULE")
				if rrule != nil {
					rRuleVal = rrule.Value
				}
				tzid := dtstart.Params.Get("TZID")
				fmt.Printf("\nParsing object: %s %s %s %s %s, exdate: %s\n",
					uid, dtstart, dtend, summary, tzid, exdateVal)
				if rRuleVal != "" {
					fmt.Printf("RRULE: %s\n", rRuleVal)
				}

				location, err := time.LoadLocation(tzid)
				if err != nil {
					return nil, fmt.Errorf("error loading location with tzid: %s: %w", tzid, err)
				}

				startTime, err := parseDate(dtstart.Value, location)
				if err != nil {
					return nil, fmt.Errorf("error parsing start date: %w", err)
				}
				endTime, err := parseDate(dtend.Value, location)
				if err != nil {
					return nil, fmt.Errorf("error parsing end date: %w", err)
				}

				if rRuleVal != "" {
					recEvent := Event{uid.Value, startTime, endTime, summary.Value}
					eventInstances, err := getRecurrenceEvents(recEvent, reqStart, reqEnd, rRuleVal, exdateVal, location)
					if err != nil {
						return nil, fmt.Errorf("error parsing events based on rrule: %s, err: %w", rRuleVal, err)
					}
					events = append(events, eventInstances...)
				} else {
					event := Event{uid.Value, startTime, endTime, summary.Value}
					events = append(events, event)
				}
			}
		}
	}

	// Sort events
	sort.Slice(events, func(i, j int) bool {
		if events[i].Start.Equal(events[j].Start) {
			return events[i].End.Before(events[j].End)
		} else {
			return events[i].Start.Before(events[j].Start)
		}
	})
	for _, event := range events {
		fmt.Printf("Event: %s %s %s %s\n",
			event.Uid, event.Start.Format(time.DateTime), event.End.Format(time.DateTime), event.Summary)
	}
	return events, nil
}

func getRecurrenceEvents(event Event, from, to time.Time, rRuleVal string, exdateStr string, tz *time.Location) ([]Event, error) {
	// Calculate event start and end time difference
	eventDuration := event.End.Sub(event.Start)
	// Parse the recurrence rule
	rule, err := rrule.StrToRRule(rRuleVal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rrule: %w", err)
	}
	rule.DTStart(event.Start)

	// Generate the dates between from and to
	dates := rule.Between(event.Start, to, true)

	// Parse the exdate string
	var exDates []time.Time
	if exdateStr != "" {
		exdateVals := strings.Split(exdateStr, ",")
		for _, exdate := range exdateVals {
			t, err := parseDate(exdate, tz)
			if err != nil {
				return nil, fmt.Errorf("failed to parse exdate: %w", err)
			}
			exDates = append(exDates, t)
		}
	}

	// Create a map for faster lookup of exDate
	exDate := make(map[time.Time]bool)
	for _, exdate := range exDates {
		exDate[exdate] = true
	}

	// Create events, filtering out the exdates and events outside the from-to range
	var events []Event
	for _, eventFrom := range dates {
		eventTo := eventFrom.Add(eventDuration)
		isOutsideRange := eventFrom.After(to) || eventTo.Before(from)
		if !exDate[eventFrom] && !isOutsideRange {
			events = append(events, Event{event.Uid, eventFrom, eventTo, event.Summary})
		}
	}

	fmt.Printf("Recurrence events: %v\n", events)
	return events, nil
}

// parseDate takes a date string and a timezone location,
// and returns the parsed date as a time.Time value in the base timezone provided in env variable.
// If the date string cannot be parsed, it returns an error.
//
// The date string should be in the format "20060102T150405" or "20060102".
// The timezone location should be a valid IANA Time Zone database name
// (e.g., "America/New_York"). If the location is nil, the function uses UTC.
//
// Example:
//
//	t, err := parseDate("20220412T123000", time.LoadLocation("America/New_York"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(t)
//
// This will print the time corresponding to April 12, 2022, 12:30:00 in the New York timezone.
func parseDate(d string, tz *time.Location) (time.Time, error) {
	parsed, err := time.ParseInLocation("20060102T150405", d, tz)
	if err != nil {
		parsed, err = time.ParseInLocation("20060102", d, tz)
		if err != nil {
			return time.Time{}, err
		}
	}
	parsed = parsed.In(config.baseTimezone)
	return parsed, nil
}

func buildQuery(from DateOffset, to DateOffset) caldav.CalendarQuery {
	return caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
				Props: []string{
					"SUMMARY",
					"UID",
					"DTSTART",
					"DTEND",
					"RRULE",
					"EXDATE",
				},
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name:  "VEVENT",
				Start: time.Now().AddDate(from.Years, from.Months, from.Days),
				End:   time.Now().AddDate(to.Years, to.Months, to.Days),
			}},
		},
	}
}
