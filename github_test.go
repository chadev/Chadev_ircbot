// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestValidateURL(t *testing.T) {
	if !validateURL("https://github.com/chadev/Chadev_ircbot") {
		t.Error("failed to validate proper URL: https://github.com/chadev/Chadev_ircbot")
	}
}

func TestGetGitHubURL(t *testing.T) {
	URL, err := getGitHubURL("Chadev_ircbot")
	if err != nil {
		t.Errorf("failed fetching GitHub repo URL: %s\n", err.Error())
	}

	t.Logf("returned URL: %s\n", URL)
	if !validateURL(URL) {
		t.Error("the URL came back as invalid")
	}
}

func TestGetIssueURL(t *testing.T) {
	URL, err := getIssueURL("Chadev_ircbot")
	if err != nil {
		t.Errorf("failed fetching GitHub Issue URL: %s\n", err.Error())
	}

	t.Logf("returned URL: %s\n", URL)
	if !validateURL(URL) {
		t.Error("the URL came back as invalid")
	}
}

func TestGetIssueIDURL(t *testing.T) {
	URL, err := getIssueURL("Chadev_ircbot")
	if err != nil {
		t.Errorf("failed fetching GitHub Issue URL: %s\n", err.Error())
	}

	t.Logf("returned URL: %s\n", URL)
	if !validateURL(URL) {
		t.Error("the issue queue URL came back as invalid")
	}

	URL, err = getIssueIDURL(URL, "#1")
	if err != nil {
		t.Errorf("failed fetching GitHub Issue URL: %s\n", err.Error())
	}

	t.Logf("returned URL: %s\n", URL)
	if !validateURL(URL) {
		t.Errorf("the issue URL came back as invalid")
	}
}
