// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/danryan/hal"
)

var musicHandler = hear(`give me some music`, "give me some music", "Returns a list of music playlists popular with the community", func(res *hal.Response) error {
	pl := []string{
		"http://open.spotify.com/user/juzten/playlist/0yCFUrwFvi4lu19DfkMcuH",
	}

	list := strings.Join(pl, ",")

	return res.Send(fmt.Sprintf("%s: here try out these awesome tunes %v", res.UserName(), list))
})
