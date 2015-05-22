// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package meetup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Base URL for the Meetup API
// TODO shouldn't contain more than just the base path, the rest should be
// built programatically to allow for calls other than just for events.
const baseURL = "https://api.meetup.com/2/events?&sign=true&photo-host=secure&group_urlname=%s&page=20&key=%s"

// Meetup contains the ruturn value from the Meetup API.
type Meetup struct {
	// Results contains various data contained in the results object.
	Results []Results `json:"results"`
	// Meta contains various data contained in the metadata object.
	Meta Meta `json:"meta"`
}

// Results contains details for each event in the return value.
type Results struct {
	// Venue contains various data about the event's venue
	Venue Venue `json:"venue"`
	// EventURL is the URL to the meetup page for the event
	EventURL string `json:"event_url"`
	// Name in the event name from Meetup
	Name string `json:"name"`
	// Time is the start time of the event in nanoseconds since Epoch
	Time int64 `json:"time"`
	// HeadCount is number of people attended
	HeadCount int `json:"headcount"`
	// YesRSVP is the number of people that responded yes to their rsvp
	YesRSVP int `json:"yes_rsvp_count"`
	// MaybeRSVP is the nubmer of people that responded maybe to their rsvp
	MaybeRSVP int `json:"maybe_rsvp_count"`
	// RSVPLimit is the maximum number of people that can attend this event
	RSVPLimit int `json:"rsvp_limit"`
}

// Venue contains details about the venue for each event.
type Venue struct {
	// Name is the name of the location the event is being held in
	Name string `json:"name"`
}

// Meta cotains additional data returned by the Meetup API
type Meta struct {
	// TotolCount is the number of events planned for a group
	TotalCount int `json:"total_count"`
}

func (e *Results) parseDateTime() time.Time {
	// set the timezone, otherwise UTC will be used
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		// couldn't set timezone, falling back to UTC
		loc, _ = time.LoadLocation("UTC")
	}

	// parse the unix timestamp, and apply our timezone
	t := time.Unix(0, e.Time*int64(time.Millisecond)).In(loc)

	return t
}

func (e *Results) parseRSVPs() string {
	var output []string

	if e.HeadCount > 0 {
		output = append(output, fmt.Sprintf("Expected headcount %d", e.HeadCount))
	}

	if e.YesRSVP > 0 {
		output = append(output, fmt.Sprintf("Confirmed RSVPs %d", e.YesRSVP))
	}

	if e.MaybeRSVP > 0 {
		output = append(output, fmt.Sprintf("Maybes %d", e.MaybeRSVP))
	}

	if e.RSVPLimit > 0 {
		output = append(output, fmt.Sprintf("Note this event is capped at %d attendiees", e.RSVPLimit))
	}

	return strings.Join(output, ", ")
}

func formatDateTime(dt time.Time) string {
	today := isToday(dt)

	var fTime string
	if today {
		fTime = dt.Format("3:04 pm")
	} else {
		fTime = dt.Format("01/02 3:04 pm")
	}

	return fTime
}

func isToday(dt time.Time) bool {
	now := time.Now().Format("2006-01-02")
	event := dt.Format("2006-01-02")

	if event != now {
		return false
	}

	return true
}

func (m *Meetup) string() string {
	dt := m.Results[0].parseDateTime()
	fTime := formatDateTime(dt)
	lunchDay := isToday(dt)

	if !lunchDay {
		return fmt.Sprintf("The next meetup is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
			m.Results[0].Name,
			m.Results[0].Venue.Name,
			fTime,
			m.Results[0].EventURL)
	}

	return fmt.Sprintf("The meetup today is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
		m.Results[0].Name,
		m.Results[0].Venue.Name,
		fTime,
		m.Results[0].EventURL)
}

// GetNextMeetup returns details for a groups next event, group is the group
// name as seen in the meetup URL.  Returns an empty string if the group has no
// upcoming events.
func GetNextMeetup(group string) (string, error) {
	Events, err := getMeetupResponce(group)
	if err != nil {
		return "", err
	}

	if Events.Meta.TotalCount == 0 {
		return "", nil
	}

	return Events.string(), nil
}

// GetMeetupRSVP returns RSVP details for a groups next event, group is the group
// name as seen in the meetup URL.  Returns an empty string if the group has no
// upcoming events.
func GetMeetupRSVP(group string) (string, error) {
	Events, err := getMeetupResponce(group)
	if err != nil {
		return "", err
	}

	if Events.Meta.TotalCount == 0 {
		return "", nil
	}

	return Events.Results[0].parseRSVPs(), nil
}

func getMeetupResponce(g string) (Meetup, error) {
	var e Meetup
	u := fmt.Sprintf(baseURL, g, os.Getenv("CHADEV_MEETUP"))
	r, err := http.Get(u)
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return e, errors.New("failed to fetch details from Meetup API")
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return e, err
	}

	err = json.Unmarshal(b, &e)
	if err != nil {
		return e, err
	}

	return e, nil
}
