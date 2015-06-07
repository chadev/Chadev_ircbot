// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danryan/hal"
)

var karmaHandler = hear(`(.+)(\+\+|\-\-)`, "(nickname)++ or (nickname)--", "Increases or Decreases a persons karma", func(res *hal.Response) error {
	nick := res.Match[1]
	sign := res.Match[2]

	var karma int
	key := strings.ToUpper(nick)
	val, _ := res.Robot.Store.Get("Karma::" + key)
	num, _ := strconv.ParseInt(string(val), 10, 0)
	karma = int(num)

	if sign == "++" {
		karma++
	} else if sign == "--" {
		karma--
	} else {
		hal.Logger.Errorf("invalid sign '%v' given", sign)
		err := res.Reply("that is not a valid option, please try again with either ++ or --")
		return err
	}

	s := strconv.Itoa(karma)
	err := res.Robot.Store.Set("Karma::"+key, []byte(s))
	if err != nil {
		hal.Logger.Errorf("could not update the user's karma: %v", err)
		err = res.Reply("I was unable to update the users karma, please try again later")
		return err
	}

	err = res.Send(fmt.Sprintf("Karma for %s has been updated", nick))
	return err
})

var karmaStatsHandler = hear(`karma stats(.*)`, "karma stats (username)", "Shows the current karma stats for a user, defaults to the user that ran the cammand", func(res *hal.Response) error {
	nick := strings.TrimSpace(res.Match[1])
	if nick == "" {
		nick = res.UserName()
	}
	key := strings.ToUpper(nick)

	val, err := res.Robot.Store.Get("Karma::" + key)
	if err != nil {
		hal.Logger.Infof("user %s has no karma", nick)
		err = res.Reply(fmt.Sprintf("User %s currently has no karma", nick))
		return err
	}

	karma, _ := strconv.Atoi(string(val))

	err = res.Send(fmt.Sprintf("User %s currently has %d karma", nick, karma))
	return err
})
