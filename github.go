// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/danryan/hal"
)

var sourceHandler = hear(`source (.+)`, func(res *hal.Response) error {
	URL, err := getGitHubURL(res.Match[1])
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("unable to get GitHub URL: %v\n", err))
		return res.Send(fmt.Sprintf("Fetching URL for %s failed, possibly misspelled?", res.Match[1]))
	}

	return res.Send(URL)
})

var issueHandler = hear(`issue (.+)`, func(res *hal.Response) error {
	URL, err := getIssueURL(res.Match[1])
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("unable to get issue URL: %v\n", err))
		return res.Send(fmt.Sprintf("Fetching issue queue URL for %s failed, possibly misspelled?", res.Match[1]))
	}

	return res.Send(URL)
})

func getGitHubURL(s string) (string, error) {
	if s == "" {
		s = "Chadev_ircbot"
	}

	// build the GitHub URL
	URL := fmt.Sprintf("https://github.com/chadev/%s", url.QueryEscape(s))

	if !validateURL(URL) {
		return "", errors.New("unable to get GitHub URL: no repo with URL: " + URL)
	}

	return URL, nil
}

func getIssueURL(s string) (string, error) {
	if s == "" {
		s = "Chadev_ircbot"
	}

	// build the URL
	URL := fmt.Sprintf("https://github.com/chadev/%s/issues", url.QueryEscape(s))

	if !validateURL(URL) {
		return "", errors.New("unable to get GitHub URL: no repo with URL: " + URL)
	}

	return URL, nil
}

func validateURL(u string) bool {
	// check if the URL is valid
	resp, err := http.Get(u)
	if err != nil {
		return false
	}

	if resp.StatusCode != 200 {
		return false
	}

	return true
}
