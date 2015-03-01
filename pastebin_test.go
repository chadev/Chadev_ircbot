// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"testing"
)

func TestPastetoPastebin(t *testing.T) {
	if os.Getenv("CHADEV_PASTEBIN") == "" {
		t.Skip("missing pastebin API key skipping tist")
	}

	u := uploadHelpMsg("This is a test post")
	t.Logf("URL: %s", u)

	if strings.Contains("github.com", u) {
		t.Error("failed to paste to pastebin, fallback url returned")
	}

	if !validateURL(u) {
		t.Error("URL came back as invalide")
	}
}
