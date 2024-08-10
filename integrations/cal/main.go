package cal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	webdav "github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/joho/godotenv"
)

type Event struct {
	Uid     string
	Start   time.Time
	End     time.Time
	Summary string
}

type DateOffset struct {
	Years  int
	Months int
	Days   int
}

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
	if (from == DateOffset{}) {
		from = DateOffset{Years: 0, Months: 0, Days: -7}
	}
	if (to == DateOffset{}) {
		to = DateOffset{Years: 0, Months: 0, Days: 7}
	}

	println("Getting family calendar events from: ", from.String(), " to: ", to.String())

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	httpClient := &http.Client{}
	username := os.ExpandEnv("$CAL_USERNAME")
	password := os.ExpandEnv("$CAL_PASSWORD")
	calUrl := os.ExpandEnv("$CAL_URL")

	fmt.Printf("Connecting to %s with %s\n", calUrl, username)
	authorizedClient := webdav.HTTPClientWithBasicAuth(httpClient, username, password)
	calDavClient, err := caldav.NewClient(authorizedClient, calUrl)
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
	fmt.Println("Querying calendars with the name 'Family'")
	calQuery := buildQuery(from, to)
	for _, cal := range calendars {
		if strings.Contains(cal.Name, "Family") {
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
				tzid := dtstart.Params.Get("TZID")
				fmt.Printf("Parsing object: %s %s %s %s %s\n",
					uid, dtstart, dtend, summary, tzid)

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

				event := Event{uid.Value, startTime, endTime, summary.Value}
				events = append(events, event)
			}
		}
	}
	for _, event := range events {
		fmt.Printf("Event: %s %s %s %s\n",
			event.Uid, event.Start.Format(time.DateTime), event.End.Format(time.DateTime), event.Summary)
	}
	return events, nil
}

func parseDate(d string, location *time.Location) (time.Time, error) {
	parsed, err := time.ParseInLocation("20060102T150405", d, location)
	if err != nil {
		parsed, err = time.ParseInLocation("20060102", d, location)
		if err != nil {
			return time.Time{}, err
		}
	}
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
					"DURATION",
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
