package main

import (
	"os"
	"regexp"
	"testing"
)

func TestGetCalendarEvents(t *testing.T) {
	if !checkOAuth() {
		t.Skip("missing one or more OAuth2 arguments, skipping test")
	}

	resp, err := getCalendarEvents()
	if err != nil {
		t.Errorf("failed to fetch events from the calendar: %v\n", err)
	}

	t.Logf("%#v\n", resp)
	match, err := regexp.MatchString("Next 7 events: (.+)", resp)
	if err != nil {
		t.Errorf("error when checking return: %v\n", err)
	}
	if !match {
		t.Error("the returned value is not what was expected")
	}
}

func TestGetNextEvent(t *testing.T) {
	if !checkOAuth() {
		t.Skip("missing one or more OAuth2 arguments, skipping test")
	}

	resp, err := getNextEvent()
	if err != nil {
		t.Errorf("failed to fetch event from the calendar: %v\n", err)
	}

	t.Logf("%#v\n", resp)
	match, err := regexp.MatchString("Next event: ", resp)
	if err != nil {
		t.Errorf("error when checking return: %v\n", err)
	}
	if !match {
		t.Errorf("the returned value is not what was expected")
	}
}

func checkOAuth() bool {
	clientID := os.Getenv("CHADEV_ID")
	clientSecret := os.Getenv("CHADEV_SECRET")
	refreshToken := os.Getenv("CHADEV_TOKEN")

	// if any of the OAuth2 arguments are missing skip the test
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return false
	}
	return true
}
