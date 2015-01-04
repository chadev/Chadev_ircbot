// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/danryan/hal"
)

var whatAreHandler = hear(`what (are|is) (.*)`, "what [are|is] (query)", "Has the bot search for something", func(res *hal.Response) error {
	query := res.Match[2]
	URL, err := getSearchURL(query)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("%s, sorry I wasn't able to search for that!", res.UserName())
	}

	return res.Send(fmt.Sprintf("%s, here is the search results: %s", res.UserName(), URL))
})

func getSearchURL(q string) (string, error) {
	URL := fmt.Sprintf("https://google.com/search?q=%s", url.QueryEscape(q))

	if !validateURL(URL) {
		return "", errors.New("unable to get GitHub URL: no repo with URL: " + URL)
	}

	return URL, nil
}
