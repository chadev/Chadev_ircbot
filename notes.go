// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/danryan/hal"
)

var noteStoreHandler = hear(`remember (\w+): (.+)`, "remeber (key): (value)", "Causes the bot to read back a stored note", func(res *hal.Response) error {
	key := strings.ToUpper(res.Match[1])
	msg := res.Match[2]

	err := res.Robot.Store.Set(key, []byte(msg))
	if err != nil {
		return res.Send("Sorry I can't remember that")
	}

	return res.Send("Got it!")
})

var noteGetHandler = hear(`recall (\w+)`, "recall (key)", "Tells Ash to remember something", func(res *hal.Response) error {
	key := strings.ToUpper(res.Match[1])
	val, err := res.Robot.Store.Get(key)
	if err != nil {
		return res.Send(fmt.Sprintf("I have no memery of %s", res.Match[1]))
	}

	return res.Send(fmt.Sprintf("%s, here is what I recall for that: %s", res.UserName(), val))
})

var noteRemoveHandler = hear(`forget (\w+)`, "forget (key)", "Tells Ash to forget something", func(res *hal.Response) error {
	key := strings.ToUpper(res.Match[1])
	_, err := res.Robot.Store.Get(key)
	if err != nil {
		return res.Send(fmt.Sprintf("I have no memery of %s", res.Match[1]))
	}

	err = res.Robot.Store.Delete(key)
	if err != nil {
		return res.Send(fmt.Sprintf("I seem to be unable to forget about %s", res.Match[1]))
	}

	return res.Send("I have forgotten it.")
})
