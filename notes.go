// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/danryan/hal"
)

var noteStoreHandler = hear(`remember ([^key].+): (.+)`, func(res *hal.Response) error {
	key := strings.ToUpper(res.Match[1])
	msg := res.Match[2]

	err := res.Robot.Store.Set(key, []byte(msg))
	if err != nil {
		return res.Send("Sorry I can't remember that")
	}

	return res.Send("Got it!")
})

var noteGetHandler = hear(`recall (.+)`, func(res *hal.Response) error {
	key := strings.ToUpper(res.Match[1])
	val, err := res.Robot.Store.Get(key)
	if err != nil {
		return res.Send(fmt.Sprintf("I have no memery of %s", res.Match[1]))
	}

	return res.Send(fmt.Sprintf("%s, here is what I recall for that: %s", res.UserName(), val))
})
