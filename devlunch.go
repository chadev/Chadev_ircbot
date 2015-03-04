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
	"strings"
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

var talkHandler = hear(`tell me about the next talk\b`, "tell me about the next talk", "Returns details on the next Chadev Lunch Talk", func(res *hal.Response) error {
	msg, err := getTalkDetails(false)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("Sorry I was unable to get details on the next dev lunch.  Please check https://meetup.com/chadevs")
	}

	return res.Send(msg)
})

var addTalkHandler = hear(`devlunch url ([a-z0-9-\s]*)(http(s)?://.+)`, "devlunch url (date) (url)", "Set live stream url for dev lunch talks", func(res *hal.Response) error {
	var d, u string
	var date time.Time

	// grab the arguments
	d = strings.TrimSpace(res.Match[1])
	u = res.Match[2]

	// if d is empty or "today" use todays date
	if d == "" || d == "today" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", d)
		if err != nil {
			// could not parse the given date, fallback to today
			hal.Logger.Error(err)
			date = time.Now()
		}
	}

	hal.Logger.Info(fmt.Sprintf("parsed date: %v", date.Format("2006-01-02")))
	if !validateURL(u) {
		return res.Send(fmt.Sprintf("%s is not a valid URL", u))
	}

	b, err := json.Marshal(DevTalk{Date: date.Format("2006-01-02"), URL: u})
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("I have failed you, I was unable to JSON")
	}

	err = res.Robot.Store.Set("devtalk", b)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("I couldn't store the live stream details")
	}

	return res.Send("Dev Talk live stream details stored")
})

var devTalkLinkHandler = hear(`link to devlunch`, "link to devlunch", "Returns the link to the dev lunch live stream", func(res *hal.Response) error {
	// check if today is Thursday
	t := time.Now()
	if t.Weekday().String() != "thursday" {
		return res.Send("Sorry today is not dev lunch day.")
	}

	// check if there is a url stored, and if the stored url is current
	b, err := res.Robot.Store.Get("devtalk")
	if err != nil || b == nil {
		hal.Logger.Error(err)
		return res.Send("Sorry, I don't have a URL for today's live stream.  You can check if it is posted to the Meeup page at http://www.meetup.com/chadevs/ or our Google+ page at https://plus.google.com/b/103401260409601780643/103401260409601780643/posts")
	}

	var talk DevTalk
	err = json.Unmarshal(b, &talk)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("Sorry, I don't have a URL for today's live stream.  You can check if it is posted to the Meeup page at http://www.meetup.com/chadevs/ or our Google+ page at https://plus.google.com/b/103401260409601780643/103401260409601780643/posts")
	}

	if talk.Date != t.Format("2006-01-02") {
		return res.Send("Sorry, I don't have a URL for today's live stream.  You can check if it is posted to the Meeup page at http://www.meetup.com/chadevs/ or our Google+ page at https://plus.google.com/b/103401260409601780643/103401260409601780643/posts")
	}

	return res.Send(fmt.Sprintf("You can access the live stream for the talk here %s", talk.URL))
})

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
