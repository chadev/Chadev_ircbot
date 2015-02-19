// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/danryan/hal"
)

var lunchHandler = hear(`is today (devlunch|dev lunch) day\b`, "is today devlunch day", "Tells if today is lunch day, and what the talk is", func(res *hal.Response) error {
	d := time.Now().Weekday().String()
	if d != "Thursday" {
		msg, err := getTalkDetails(false)
		if err != nil {
			hal.Logger.Error(err)
			return res.Send("Sorry I was unable to get details on the next dev lunch.  Please check https://meetup.com/chadevs")
		}

		return res.Send(fmt.Sprintf("No, sorry!  %s", msg))
	}

	msg, err := getTalkDetails(true)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("Sorry I was unable to get details on the next dev lunch.  Please check https://meetup.com/chadevs")
	}

	return res.Send(fmt.Sprintf("Yes!  %s", msg))
})

var talkHandler = hear(`tell me about the next talk\b`, "what is the next talk", "Returns details on the next Chadev Lunch Talk", func(res *hal.Response) error {
	msg, err := getTalkDetails(false)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("Sorry I was unable to get details on the next dev lunch.  Please check https://meetup.com/chadevs")
	}

	return res.Send(msg)
})

type Meetup struct {
	Results []MeetupEvents `json:"results"`
}

type MeetupEvents struct {
	Venue    MeetupVenue `json:"venue"`
	EventURL string      `json:"event_url"`
	Name     string      `json:"name"`
	Time     int64       `json:"time"`
}

type MeetupVenue struct {
	Name string `json:"name"`
}

func (m *Meetup) string(lunchDay bool) string {
	if !lunchDay {
		return fmt.Sprintf("The next talk is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
			m.Results[0].Name,
			m.Results[0].Venue.Name,
			m.Results[0].parseDateTime(false),
			m.Results[0].EventURL)
	}

	return fmt.Sprintf("The talk today is \"%s\", you can join us at %s on %s.  If you plan to come please make sure you have RSVPed at %s",
		m.Results[0].Name,
		m.Results[0].Venue.Name,
		m.Results[0].parseDateTime(true),
		m.Results[0].EventURL)
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

func getTalkDetails(lunchDay bool) (string, error) {
	URL := fmt.Sprintf("https://api.meetup.com/2/events?&sign=true&photo-host=secure&group_urlname=chadevs&page=20&key=%s", os.Getenv("CHADEV_MEETUP"))
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var Events Meetup
	err = json.Unmarshal(body, &Events)
	if err != nil {
		return "", err
	}

	return Events.string(lunchDay), nil
}
