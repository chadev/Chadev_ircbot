// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestGetMeetupEvents(t *testing.T) {
	_, err := getTalkDetails()
	if err != nil {
		t.Error(err)
	}
}
