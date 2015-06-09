// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/danryan/hal"
)

type DevKarma struct {
	Name  string `json:"name"`
	Karma int    `json:"karma"`
}

var karmaHandler = hear(`(.+)(\+\+|\-\-)`, "(nickname)++ or (nickname)--", "Increases or Decreases a persons karma", func(res *hal.Response) error {
	nick := res.Match[1]
	sign := res.Match[2]

	var dev DevKarma
	key := strings.ToUpper(nick)
	val, _ := res.Robot.Store.Get("Karma::" + key)
	json.Unmarshal(val, &dev)

	// check to see if we got something back
	if dev.Name == "" {
		dev.Name = nick
	}

	if sign == "++" {
		dev.Karma++
	} else if sign == "--" {
		dev.Karma--
	} else {
		hal.Logger.Errorf("invalid sign '%v' given", sign)
		err := res.Reply("that is not a valid option, please try again with either ++ or --")
		return err
	}

	val, _ = json.Marshal(dev)
	err := res.Robot.Store.Set("Karma::"+key, val)
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

	var dev DevKarma
	key := strings.ToUpper(nick)
	val, err := res.Robot.Store.Get("Karma::" + key)
	if err != nil {
		hal.Logger.Infof("user %s has no karma", nick)
		err = res.Reply(fmt.Sprintf("User %s currently has no karma", nick))
		return err
	}
	json.Unmarshal(val, &dev)

	if dev.Name == "" {
		hal.Logger.Infof("user %s has no karma", nick)
		err = res.Reply(fmt.Sprintf("User %s currently has no karma", nick))
		return err
	}

	err = res.Send(fmt.Sprintf("User %s currently has %d karma", nick, dev.Karma))
	return err
})
