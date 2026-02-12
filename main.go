package main

import (
	fmt "fmt"
	"os"
	"slices"
	"time"

	"github.com/google/uuid"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	year := 2026
	holidays := holidaysInYear(year)

	// Non-official holidays that City of Steinbach doesn't work.
	holidays = append(holidays,
		// idk what this one is
		time.Date(year, time.August, 3, 0, 0, 0, 0, time.UTC),
		// Christmas Eve
		time.Date(year, time.December, 24, 0, 0, 0, 0, time.UTC),
	)

	garbageDays := garbageDaysInYear(year, holidays)

	fmt.Println("Garbage days:")
	printDays(garbageDays)
	//timezone := time.FixedZone("Steinbach", -6*60*60)
	events := []vEvent{}
	//events := []vEvent{
	//	{
	//		dateOnly:    true,
	//		uid:         fmt.Sprintf("%s@losers.fyi", uuid.New().String()),
	//		start:       time.Date(year, time.February, 14, 0, 0, 0, 0, timezone),
	//		end:         time.Date(year, time.February, 14, 0, 0, 0, 0, timezone),
	//		created:     time.Now(),
	//		description: "Garbage and Recycling Day",
	//		summary:     "Garbage and Recycling Day",
	//		alarms: []vAlarm{{
	//			action:  Display,
	//			trigger: "-P1D",
	//		}},
	//	},
	//}

	alarm := vAlarm{
		action:      Display,
		trigger:     "-P1D",
		description: "Garbage day tomorrow",
	}
	for _, gDay := range garbageDays {
		events = append(events, vEvent{
			dateOnly:    true,
			uid:         fmt.Sprintf("%s@losers.fyi", uuid.New().String()),
			start:       gDay,
			end:         gDay,
			created:     time.Now(),
			description: "Garbage and Recycling Day",
			summary:     "Garbage and Recycling Day",
			alarms:      []vAlarm{alarm},
		})
	}

	cal := newVCalender(events)

	fmt.Println(cal)

	err := os.WriteFile("/media/shuna/static/public/cal/garbage-days.ics", []byte(cal.String()), 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func timeArrayContains(source []time.Time, search time.Time) bool {
	for _, day := range source {
		if search.Compare(day) == 0 {
			return true
		}
	}
	return false
}

func printDays(days []time.Time) {
	prevMonth := time.January
	for _, day := range days {
		if prevMonth != day.Month() {
			fmt.Println()
			prevMonth = day.Month()
		}
		fmt.Printf("%s ", day.Format("2006-01-02"))
	}
	fmt.Println()
}

func holidaysInYear(year int) []time.Time {
	holidays := []time.Time{}
	first := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	holidays = append(holidays, first)

	// 3rd Monday in feb
	holidays = append(holidays, dayOfWeekAway(
		// subtract a day in case the 1st is a Monday
		time.Date(year, time.February, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		3, time.Monday, After),
	)

	// Good Friday - Friday before Easter Sunday (Easter Sunday is variable, and hard to determine year-by-year)
	holidays = append(holidays, time.Date(year, time.April, 3, 0, 0, 0, 0, time.UTC))

	// Victoria Day - Monday before May 25th
	holidays = append(holidays, dayOfWeekAway(
		time.Date(year, time.May, 25, 0, 0, 0, 0, time.UTC),
		1, time.Monday, Before),
	)

	// Canada Day
	holidays = append(holidays, time.Date(year, time.July, 1, 0, 0, 0, 0, time.UTC))

	// Labour Day - First Monday in September
	holidays = append(holidays, dayOfWeekAway(
		// subtract a day in case the 1st is a Monday
		time.Date(year, time.September, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		3, time.Monday, After),
	)

	// National Day for Truth and Reconciliation
	holidays = append(holidays, time.Date(year, time.September, 30, 0, 0, 0, 0, time.UTC))

	// Thanksgiving - Second Monday in October
	holidays = append(holidays, dayOfWeekAway(
		time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		2, time.Monday, After),
	)

	// Remembrance day
	holidays = append(holidays, time.Date(year, time.November, 11, 0, 0, 0, 0, time.UTC))

	// Christmas day
	holidays = append(holidays, time.Date(year, time.December, 25, 0, 0, 0, 0, time.UTC))

	// Boxing day
	holidays = append(holidays, time.Date(year, time.December, 26, 0, 0, 0, 0, time.UTC))

	return holidays
}

func businessDaysInYear(year int, holidays []time.Time) []time.Time {
	businessDays := []time.Time{}

	daysInYear := 365
	workingDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	for range daysInYear {
		workingDate = workingDate.AddDate(0, 0, 1)
		if workingDate.Weekday() == time.Saturday || workingDate.Weekday() == time.Sunday {
			continue
		}

		if timeArrayContains(holidays, workingDate) {
			continue
		}

		businessDays = append(businessDays, workingDate)
	}
	return businessDays
}

func garbageDaysInYear(year int, holidays []time.Time) []time.Time {
	businessDays := businessDaysInYear(year, holidays)
	startDate := time.Date(year, time.January, 7, 0, 0, 0, 0, time.UTC)
	startIndex := slices.Index(businessDays, startDate)

	if startIndex == -1 {
		panic("No start date found")
	}

	garbageDays := []time.Time{}
	i := startIndex
	for i < len(businessDays) {
		garbageDays = append(garbageDays, businessDays[i])
		i += 5
	}
	return garbageDays
}

type Direction int

const (
	Before Direction = -1
	After  Direction = 1
)

// Returns the date, not including start date, that is `days` number of `weekday` weekdays away in the corresponding direction.
func dayOfWeekAway(date time.Time, days int, weekday time.Weekday, direction Direction) time.Time {
	counted := 0
	next := date

	for counted < days {
		next = next.AddDate(0, 0, int(direction))
		if next.Weekday() == weekday {
			counted++
		}
	}
	return next
}
