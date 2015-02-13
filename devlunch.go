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

var lunchHandler = hear(`is today (devlunch|dev lunch) day?`, "is today devlunch day", "Tells if today is lunch day, and what the talk is", func(res *hal.Response) error {
	d := time.Now().Weekday().String()
	if d != "Thursday" {
		//return res.Send("No today is not dev lunch day, sorry!")
	}

	msg, err := getTalkDetails()
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
}

type MeetupVenue struct {
	Name string `json:"name"`
}

func (m *Meetup) string() string {
	return fmt.Sprintf("Yes!  The talk today is \"%s\", you can join us at %s.  If you plan to come please make sure you have RSVPed at %s", m.Results[0].Name, m.Results[0].Venue.Name, m.Results[0].EventURL)
}

func getTalkDetails() (string, error) {
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

	return Events.string(), nil
}
