package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/danryan/hal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var eventHandler = hal.Hear(`.events`, func(res *hal.Response) error {
	events, err := getCalendarEvents()
	if err != nil {
		return res.Send("Could not fetch list of events")
	}

	return res.Send(events)
})

var baseURL = "https://www.googleapis.com/calendar/v3"

func getCalendarEvents() (string, error) {
	oauthToken := os.Getenv("CHADEV-TOKEN")
	var opts *oauth2.Options
	var err error
	if oauthToken == "" {
		opts, err = getOauth2Token()
		if err != nil {
			return "", err
		}
	}

	url := opts.AuthCodeURL("test", "offline", "auto")
	return url, nil

	URL := fmt.Sprintf("%s/calendars/%s/events", baseURL, "4qc3thgj9ocunpfist563utr6g")
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getOauth2Token() (*oauth2.Options, error) {
	clientID := os.Getenv("CHADEV-ID")
	clientSecret := os.Getenv("CHADEV-SECRET")
	opts, err := oauth2.New(
		oauth2.Client(clientID, clientSecret),
		oauth2.RedirectURL("urn:ietf:wg:oauth:2.0:oob"),
		oauth2.Scope("https://www.googleapis.com/auth/calendar.readonly"),
		google.Endpoint(),
	)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
