package main

import (
	"fmt"
	"strings"
	"time"
)

type vCalendar struct {
	events []vEvent

	// Identifier of the program used to produce the iCal
	prodid  string // required
	version string // required
}

func newVCalender(events []vEvent) *vCalendar {
	cal := vCalendar{
		version: "2.0",
		prodid:  "GoLang iCalendar - custom",
		events:  events,
	}

	return &cal
}

func (cal vCalendar) String() string {
	builder := strings.Builder{}
	builder.WriteString("BEGIN:VCALENDAR\r\n")
	builder.WriteString(fmt.Sprintf("PRODID:%s\r\n", cal.prodid))
	builder.WriteString(fmt.Sprintf("VERSION:%s\r\n", cal.version))
	for _, event := range cal.events {
		builder.WriteString(event.String())
	}
	builder.WriteString("END:VCALENDAR\r\n")
	return builder.String()
}

type vEvent struct {
	dateOnly    bool
	uid         string
	start       time.Time
	end         time.Time
	created     time.Time
	description string
	summary     string
	alarms      []vAlarm
}

func (event vEvent) add(alarm vAlarm) *vEvent {
	event.alarms = append(event.alarms, alarm)
	return &event
}

func (event vEvent) String() string {
	builder := strings.Builder{}
	builder.WriteString("BEGIN:VEVENT\r\n")
	builder.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", event.start.UTC().Format("20060102T150405Z")))
	builder.WriteString(fmt.Sprintf("UID:%s\r\n", event.uid))
	if event.dateOnly {
		builder.WriteString(fmt.Sprintf("DTSTART;VALUE=DATE:%s\r\n", event.start.Format("20060102")))
		builder.WriteString(fmt.Sprintf("DTEND;VALUE=DATE:%s\r\n", event.start.Format("20060102")))
	} else {
		builder.WriteString(fmt.Sprintf("DTSTART:%s\r\n", event.start.Format("20060102T150405")))
		builder.WriteString(fmt.Sprintf("DTEND:%s\r\n", event.start.Format("20060102T150405")))
	}
	builder.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", event.description))
	builder.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", event.summary))
	for _, alarm := range event.alarms {
		builder.WriteString(alarm.String())
	}
	builder.WriteString("END:VEVENT\r\n")
	return builder.String()
}

type AlarmActionType int

const (
	Audio AlarmActionType = iota
	Display
	Email
)

var actionTypes = []string{"AUDIO", "DISPLAY", "EMAIL"}

type vAlarm struct {
	action      AlarmActionType
	trigger     string
	description string
}

func (alarm vAlarm) String() string {
	builder := strings.Builder{}
	builder.WriteString("BEGIN:VALARM\r\n")

	builder.WriteString(fmt.Sprintf("ACTION:%s\r\n", actionTypes[alarm.action]))
	builder.WriteString(fmt.Sprintf("TRIGGER;RELATED=START:%s\r\n", alarm.trigger))

	builder.WriteString("END:VALARM\r\n")
	return builder.String()
}
