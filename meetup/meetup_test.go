// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package meetup

import (
	"os"
	"testing"
	"time"
)

func TestGetMeetupEvents(t *testing.T) {
	if os.Getenv("CHADEV_MEETUP") == "" {
		t.Skip("no meetup API key set, skipping test")
	}

	var l bool
	d := time.Now().Weekday().String()
	if d == "Thursday" {
		l = true
	}

	_, err := GetTalkDetails(l)
	if err != nil {
		t.Error(err)
	}
}
