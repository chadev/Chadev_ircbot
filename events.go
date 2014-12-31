// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/danryan/hal"
)

var eventHandler = hear(`events`, "events", "Get next 7 events from the Chadev calendar", func(res *hal.Response) error {
	events, err := getCalendarEvents()
	if err != nil {
		hal.Logger.Error("failed to call Calendar API: %v", err)
		return res.Send("Could not fetch data from Google Calendar API, please try again later")
	}

	return res.Send(events)
})

const baseURL = "https://www.googleapis.com/calendar/v3/calendars"

var accessToken AccessToken

// AccessToken contains the current oauth2 token and its expiry time.
type AccessToken struct {
	Token string
	Valid time.Time
}

func (a *AccessToken) getExpireTime(o int) {
	n := time.Now()
	d := time.Duration(o) * time.Second

	a.Valid = n.Add(d)
}

func (a *AccessToken) expiredToken() bool {
	n := time.Now()
	if a.Valid.After(n) {
		hal.Logger.Info("current oauth token is expired")
		a.Token = ""
		return true
	}

	return false
}

// Event is the top-level of the JSON return from the calendar API.
type Event struct {
	// Kind is the type of collection ("calendar#events")
	Kind string `json:"kind"`
	// Etag is the etag of the collection
	Etag string `json:"etag"`
	// Summary is the name of the calendar
	Summary string `json:"summary"`
	// Description is the description of the calendar
	Description string `json:"descritpion"`
	// Updated is the last modified time for the calendar
	Updated string `json:"updated"`
	// TimeZone is the timezone for the calendar
	TimeZone string `json:"timeZone"`
	// AccessRole is the role of the current user.
	// possible values include: "none", "freeBuzyReader", "reader", "writer", "owner"
	AccessRole string `json:"accessRole"`
	// NextPageToken is used to access the next page of this result. Omitted if no
	// further results are available, in which case nextSyncToken is provided.
	NextPageToken string `json:"nextPageToken,omitempty"`
	// Items is the list of events on the calendar.
	Items []EventItem `json:"items"`
	// NextSyncToken is used at a later point in time to retrieve only the entries
	// that have changed since this result was returned. Omitted if further results
	// are available, in which case nextPageToken is provided.
	NextSyncToken string `json:"nextSyncToken"`
}

// EventItem contains items with-in the "items" list of the JSON.
type EventItem struct {
	// Kind is the type of collection ("calendar#event")
	Kind string `json:"kind"`
	// Etag string `json:"etag"`
	Etag string `json:"etag"`
	// ID is the ID of the event
	ID string `json:"id"`
	// Status is the status of the event
	Status string `json:"status"`
	// HTMLLink is the link to the event on the calendar
	HTMLLink string `json:"htmlLink"`
	// Created is the datetime stamp that the event was created
	Created string `json:"created"`
	// Updated is the datetame stamp that the event was last modified
	Updated string `json:"updated"`
	// Summary is the name of the event
	Summary string `json:"summary"`
	// Description is the body of the event
	Description string `json:"description"`
	// Location is the address for the event
	Location string `json:"location"`
	// Creator is the user that created the event
	Creator EventCreator `json:"creator"`
	// Organizer is the organizer information
	organizer EventOrganizer `json:"organizer"`
	// Start is the start date and time
	Start EventDateTime `json:"start"`
	// End is the end date and time
	End EventDateTime `json:"end"`
	// ICalUID is the iCal ID for the event
	ICalUID string `json:"uCalUID"`
	// Sequence denotes if this is a repeating event
	Sequence int `json:"sequence"`
	// Reminders are the event reminders
	Reminders EventReminder `json:"reminders"`
}

// EventCreator contains the fields for the creator object.
type EventCreator struct {
	// Email address for the event creator
	Email string `json:"email"`
	// DisplayName is the real name of the creator
	DisplayName string `json:"displayName"`
}

// EventOrganizer contains the fields for the organizer object.
type EventOrganizer struct {
	// Email address for the event organizer
	Email string `json:"email"`
	// DisplayName is the name of the event organizer
	DisplayName string `json:"displayName"`
	// Self denotes if the organizer is the current user
	Self bool `json:self"`
}

// EventDateTime contains the fields for the "start" and "end" objects.
type EventDateTime struct {
	// Date is the start/end date, used for all day/multi day events
	Date string `json:"date,omitempty"`
	// DateTime is the start/end datetime in  RFC 3339 format
	DateTime string `json:"dateTime,omitempty"`
	// TimeZone is the Timezone for the event
	TimeZone string `json:"timeZone,omitempty"`
}

//EventReminder contains the fields for the reminders object.
type EventReminder struct {
	// UseDefault denotes if Google Calendar default reminders are used
	UseDefault bool `json:"useDefault"`
}

func getCalendarEvents() (string, error) {
	var err error

	if accessToken.Token == "" ||
		accessToken.expiredToken() {
		accessToken, err = getOauth2Token()
		if err != nil {
			return "", err
		}
	}

	URL := fmt.Sprintf("%s/4qc3thgj9ocunpfist563utr6g@group.calendar.google.com/events?access_token=%s&singleEvents=true&orderBy=startTime&timeMin=%s&maxResults=7",
		baseURL, url.QueryEscape(accessToken.Token),
		url.QueryEscape(time.Now().Format("2006-01-02T15:04:05Z")))
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var Events Event

	err = json.Unmarshal(body, &Events)
	if err != nil {
		return "", err
	}

	eventList := getEventList(Events)

	return eventList, nil
}

func getOauth2Token() (AccessToken, error) {
	type Responce struct {
		Token   string `json:"access_token"`
		Expires int    `json:"expires_in"`
	}

	clientID := os.Getenv("CHADEV_ID")
	clientSecret := os.Getenv("CHADEV_SECRET")
	refreshToken := os.Getenv("CHADEV_TOKEN")

	if clientID == "" {
		return accessToken, errors.New("client ID is undefined")
	}

	if clientSecret == "" {
		return accessToken, errors.New("client secret is undefined")
	}

	if refreshToken == "" {
		return accessToken, errors.New("oauth refresh token is undefined")
	}

	var r Responce

	body := fmt.Sprintf("client_id=%s&client_secret=%s&refresh_token=%s&grant_type=refresh_token",
		clientID, clientSecret, refreshToken)
	b := strings.NewReader(body)

	resp, err := http.Post("https://accounts.google.com/o/oauth2/token", "application/x-www-form-urlencoded", b)
	if err != nil {
		return accessToken, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accessToken, err
	}

	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return accessToken, err
	}

	accessToken.Token = r.Token
	accessToken.getExpireTime(r.Expires)

	return accessToken, nil
}

func getEventList(events Event) string {
	var output string

	output = "Next 7 events: "
	for key, event := range events.Items {
		if key == 0 {
			output += event.Summary
		} else {
			output += ", " + event.Summary
		}
		output += " ("
		if event.Start.DateTime != "" {
			output += formatDatetime(event.Start.DateTime)
		} else {
			output += formatDate(event.Start.Date, event.End.Date)
		}
		output += ")"
	}

	return output
}

func formatDatetime(s string) string {
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	return date.Format("01/02 3:04 pm")
}

func formatDate(s, e string) string {
	var o string
	if s == e {
		date, err := time.Parse("2006-01-02", s)
		if err != nil {
			return ""
		}
		o = date.Format("01/02")
	} else {
		sDate, err := time.Parse("2006-01-02", s)
		if err != nil {
			return ""
		}
		eDate, err := time.Parse("2006-01-02", e)
		if err != nil {
			return ""
		}
		o = sDate.Format("01/02") + " - " + eDate.Format("01/02")
	}

	return o
}
