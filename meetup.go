package main

import "time"

// Meetup contains the ruturn value from the Meetup API.
type Meetup struct {
	Results []MeetupEvents `json:"results"`
}

// MeetupEvents contains details for each event in the return value.
type MeetupEvents struct {
	Venue    MeetupVenue `json:"venue"`
	EventURL string      `json:"event_url"`
	Name     string      `json:"name"`
	Time     int64       `json:"time"`
}

// MeetupVenue contains details about the venue for each event.
type MeetupVenue struct {
	Name string `json:"name"`
}

// DevTalk contains the dev talk live stream details.
type DevTalk struct {
	Date, URL string
}

func (e *MeetupEvents) parseDateTime(today bool) string {
	// set the timezone, otherwise UTC will be used
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// couldn't set timezone, falling back to UTC
		loc, _ = time.LoadLocation("UTC")
	}

	// parse the unix timestamp, and apply our timezone
	t := time.Unix(0, e.Time*int64(time.Millisecond)).In(loc)

	var fT string
	if today {
		fT = t.Format("3:04 pm")
	} else {
		fT = t.Format("01/02 3:04 pm")
	}

	return fT
}
