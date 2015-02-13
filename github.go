// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/danryan/hal"
)

var sourceHandler = hear(`source(.*)`, "source", "Give the URL to the named GitHub repo", func(res *hal.Response) error {
	URL, err := getGitHubURL(strings.TrimSpace(res.Match[1]))
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("unable to get GitHub URL: %v\n", err))
		return res.Send(fmt.Sprintf("Fetching URL for %s failed, possibly misspelled?", res.Match[1]))
	}

	return res.Send(URL)
})

var issueHandler = hear(`issue(.*)`, "issue", "Give the URL to the issue queue for the named GitHub repo", func(res *hal.Response) error {
	args := make([]string, 2)
	if res.Match[1] != "" {
		res.Match[1] = strings.TrimSpace(res.Match[1]) // trim leading and tailing whitespace from the match
		// check if we have a project name and issue nubmer "projectName issueID"
		if strings.Contains(res.Match[1], " ") {
			args = strings.Split(res.Match[1], " ")
		} else {
			args[0] = res.Match[1]
			args[1] = ""
		}
	}

	URL, err := getIssueURL(args[0])
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("unable to get issue URL: %v\n", err))
		return res.Send(fmt.Sprintf("Fetching issue queue URL for %s failed, possibly misspelled?", args[0]))
	}

	if args[1] != "" {
		URL, err = getIssueIDURL(URL, strings.TrimLeft(args[1], "#"))
		if err != nil {
			hal.Logger.Error(fmt.Sprintf("unable to get issue URL: %v\n", err))
			return res.Send(fmt.Sprintf("Fetching issue URL for issue %s failed", args[1]))
		}
	}

	return res.Send(URL)
})

func getGitHubURL(s string) (string, error) {
	if s == "" {
		s = "Chadev_ircbot"
	}
	hal.Logger.Info(s)

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

func getIssueIDURL(u, i string) (string, error) {
	// build the URL
	URL := fmt.Sprintf("%s/%s", u, i)

	if !validateURL(URL) {
		return "", errors.New("unable to get issue URL: no repo or issue with URL: " + URL)
	}

	return URL, nil
}

func validateURL(u string) bool {
	// check if the URL is valid
	_, err := url.Parse(u)
	if err != nil {
		return false
	}

	return true
}
