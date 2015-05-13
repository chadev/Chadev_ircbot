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

const baseURL = "https://api.meetup.com/2/events?&sign=true&photo-host=secure&group_urlname=%s&page=20&key=%s"

// Meetup contains the ruturn value from the Meetup API.
type Meetup struct {
	Results []MeetupEvents `json:"results"`
	Meta    MeetupMeta     `json:"meta"`
}

// MeetupEvents contains details for each event in the return value.
type MeetupEvents struct {
	Venue     MeetupVenue `json:"venue"`
	EventURL  string      `json:"event_url"`
	Name      string      `json:"name"`
	Time      int64       `json:"time"`
	HeadCount int         `json:"headcount"`
	YesRSVP   int         `json:"yes_rsvp_count"`
	MaybeRSVP int         `json:"maybe_rsvp_count"`
	RSVPLimit int         `json:"rsvp_limit"`
}

// MeetupVenue contains details about the venue for each event.
type MeetupVenue struct {
	Name string `json:"name"`
}

// MeetupMeta cotains additional data returned by the Meetup API
type MeetupMeta struct {
	TotalCount int `json:"total_count"`
}

func (e *MeetupEvents) parseDateTime() time.Time {
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

func (e *MeetupEvents) parseRSVPs() string {
	var output []string

	if e.HeadCount > 0 {
		output = append(output, fmt.Sprintf("Expected headcount: %d", e.HeadCount))
	}

	if e.YesRSVP > 0 {
		output = append(output, fmt.Sprintf("Confirmed RSVPs: %d", e.YesRSVP))
	}

	if e.MaybeRSVP > 0 {
		output = append(output, fmt.Sprintf("Maybes: %d", e.MaybeRSVP))
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

func GetNextMeetup(group string) (string, error) {
	URL := fmt.Sprintf(baseURL, group, os.Getenv("CHADEV_MEETUP"))
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Failed to fetch details from Meetup API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var Events Meetup
	err = json.Unmarshal(body, &Events)
	if err != nil {
		return "", err
	}

	if Events.Meta.TotalCount == 0 {
		return "", nil
	}

	return Events.string(), nil
}

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
