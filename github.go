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

var sourceHandler = hal.Hear(listenName+` source (.+)`, func(res *hal.Response) error {
	URL, err := getGitHubURL(res.Match[1])
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("unable to get GitHub URL: %v\n", err))
		return res.Send(fmt.Sprintf("Fetching URL for %s failed, possibly misspelled?", res.Match[1]))
	}

	return res.Send(URL)
})

func getGitHubURL(s string) (string, error) {
	// build the GitHub URL
	URL := fmt.Sprintf("https://github.com/chadev/%s", url.QueryEscape(s))

	// check if the URL is valid
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("no repo with URL: " + URL)
	}

	return URL, nil
}
