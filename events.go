package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/danryan/hal"
)

var eventHandler = hal.Hear(`.events`, func(res *hal.Response) error {
	events, err := getCalendarEvents()
	if err != nil {
		return res.Send(err.Error())
	}

	return res.Send(events)
})

var (
	baseURL     = "https://www.googleapis.com/calendar/v3/calendars"
	accessToken AccessToken
)

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
		return true
	}

	return false
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

	URL := fmt.Sprintf("%s/4qc3thgj9ocunpfist563utr6g@group.calendar.google.com/events?access_token=%s",
		baseURL, url.QueryEscape(accessToken.Token))
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

func getOauth2Token() (AccessToken, error) {
	type Responce struct {
		Token   string `json:"access_token"`
		Expires int    `json:"expires_in"`
	}

	clientID := os.Getenv("CHADEV_ID")
	clientSecret := os.Getenv("CHADEV_SECRET")
	refreshToken := os.Getenv("CHADEV_TOKEN")

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
