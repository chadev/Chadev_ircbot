// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestSearchValidateURL(t *testing.T) {
	if !validateURL("https://google.com/search?q=test") {
		t.Error("faild to validate proper URL: https://google.com/search?q=test")
	}
}

func TestGetSearchURL(t *testing.T) {
	URL, err := getSearchURL("chadev")
	if err != nil {
		t.Error("failed getting search URL: " + err.Error())
	}

	t.Log("URL: " + URL)
	if !validateURL(URL) {
		t.Error("URL came back as invalid")
	}
}
