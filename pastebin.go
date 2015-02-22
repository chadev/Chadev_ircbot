// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/danryan/hal"
)

const githubReadme = "https://github.com/chadev/Chadev_ircbot/blob/master/README.md#usage"

func uploadHelpMsg(msg string) string {
	resp, err := http.PostForm("http://pastebin.com/api/api_post.php",
		url.Values{"api_dev_key": {os.Getenv("CHADEV_PASTEBIN")},
			"api_paste_private":     {"0"},
			"api_paste_name":        {"Ash Usage"},
			"api_paste_expire_date": {"10M"},
			"api_option":            {"paste"},
			"api_paste_code":        {msg}})
	if err != nil {
		hal.Logger.Error(err)
		return githubReadme
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hal.Logger.Error(err)
		return githubReadme
	}

	url := string(body)

	if strings.Contains("Bad API request", url) {
		hal.Logger.Error(url)
		return githubReadme
	}

	return url
}
