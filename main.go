package main

import (
	"flag"
	fmt "fmt"
	"os"
	"slices"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	outputPath := flag.String("o", "/media/shuna/static/public/cal/", "Output path")
	year := flag.Int("y", 2026, "Year")
	flag.Parse()
	holidayEvents := holidayEventsInYear(*year, []vAlarm{})

	// Non-official holidays that City of Steinbach doesn't work.
	holidayEvents = append(holidayEvents,
		// Civic Holiday (august long weekend)
		*newVEvent(dayOfWeekAway(
			// subtract a day in case the 1st is a Monday
			time.Date(*year, time.August, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
			1, time.Monday, After), []vAlarm{}, "Civic Holiday"),
		//*newVEvent(time.Date(year, time.August, 3, 0, 0, 0, 0, time.UTC), []vAlarm{}, ""),
		// Christmas Eve
		*newVEvent(time.Date(*year, time.December, 24, 0, 0, 0, 0, time.UTC), []vAlarm{}, "Christmas Eve"),
	)
	holidays := holidaysInYear(holidayEvents)

	garbageDays := garbageDaysInYear(*year, holidays)

	alarm := vAlarm{
		action:      Display,
		trigger:     "-P1D",
		description: "Garbage day tomorrow",
	}

	garbageDayEvents := newVCalender(createEventsFromTime(garbageDays, []vAlarm{alarm}, "Garbage and Recycling Day"))
	writeCalendarToFile(garbageDayEvents, *outputPath+"/garbage-days.ics")
	writeCalendarToFile(newVCalender(holidayEvents), *outputPath+"/holidays.ics")
}

func createEventsFromTime(t []time.Time, alarms []vAlarm, name string) []vEvent {
	var events []vEvent
	for _, gDay := range t {
		events = append(events, *newVEvent(gDay, alarms, name))
	}
	return events
}

func writeCalendarToFile(calendar *vCalendar, path string) {
	err := os.WriteFile(path, []byte(calendar.String()), 0644)
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

func holidayEventsInYear(year int, alarms []vAlarm) []vEvent {
	holidays := []vEvent{}
	first := *newVEvent(time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC), alarms, "New Years")
	holidays = append(holidays, first)

	// 3rd Monday in feb
	holidays = append(holidays, *newVEvent(dayOfWeekAway(
		// subtract a day in case the 1st is a Monday
		time.Date(year, time.February, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		3, time.Monday, After), alarms, "Louis Riel Day"),
	)

	// Good Friday - Friday before Easter Sunday (Easter Sunday is variable, and hard to determine year-by-year)
	holidays = append(holidays, *newVEvent(
		time.Date(year, time.April, 3, 0, 0, 0, 0, time.UTC),
		alarms, "Good Friday"))

	// Victoria Day - Monday before May 25th
	holidays = append(holidays, *newVEvent(dayOfWeekAway(
		time.Date(year, time.May, 25, 0, 0, 0, 0, time.UTC),
		1, time.Monday, Before), alarms, "Victoria Day"),
	)

	// Canada Day
	holidays = append(holidays, *newVEvent(time.Date(year, time.July, 1, 0, 0, 0, 0, time.UTC),
		alarms, "Canada Day"))

	// Labour Day - First Monday in September
	holidays = append(holidays, *newVEvent(dayOfWeekAway(
		// subtract a day in case the 1st is a Monday
		time.Date(year, time.September, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		3, time.Monday, After), alarms, "Labor Day"),
	)

	// National Day for Truth and Reconciliation
	holidays = append(holidays, *newVEvent(
		time.Date(year, time.September, 30, 0, 0, 0, 0, time.UTC),
		alarms, "National Day for Truth and Reconciliation",
	))

	// Thanksgiving - Second Monday in October
	holidays = append(holidays, *newVEvent(dayOfWeekAway(
		time.Date(year, time.October, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1),
		2, time.Monday, After), alarms, "Thanksgiving Day"),
	)

	// Remembrance day
	holidays = append(holidays, *newVEvent(time.Date(year, time.November, 11, 0, 0, 0, 0, time.UTC), alarms, "Remembrance Day"))

	// Christmas day
	holidays = append(holidays, *newVEvent(time.Date(year, time.December, 25, 0, 0, 0, 0, time.UTC), alarms, "Christmas Day"))

	// Boxing day
	holidays = append(holidays, *newVEvent(time.Date(year, time.December, 26, 0, 0, 0, 0, time.UTC), alarms, "Boxing Day"))

	return holidays
}

func holidaysInYear(holidays []vEvent) []time.Time {
	dates := []time.Time{}
	for _, holiday := range holidays {
		dates = append(dates, holiday.start)
	}
	return dates
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
