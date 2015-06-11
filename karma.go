// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/danryan/hal"
)

// DevKarma contains fields for user karma information.
type DevKarma struct {
	Name  string `json:"name"`
	Karma int    `json:"karma"`
}

func (d DevKarma) String() string {
	return fmt.Sprintf("%s: %d\n", d.Name, d.Karma)
}

// ByKarma implements sort.Interface for []DevKarma based on the Karma field.
type ByKarma []DevKarma

func (a ByKarma) Len() int           { return len(a) }
func (a ByKarma) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKarma) Less(i, j int) bool { return a[i].Karma < a[j].Karma }

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

	updateDevsWithKarma(key, res)

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

var karmaRankingHandler = hear(`karma ranking(\s(top|lowest)\s\d*)?`, "karma ranking", "Shows the current karma standings, with optional top/lowest standings", func(res *hal.Response) error {
	var devs []DevKarma
	names := getDevKarmaList(res)
	hal.Logger.Debugf("Value of names: %#v", names)

	for _, v := range names {
		var dev DevKarma
		val, err := res.Robot.Store.Get("Karma::" + v)
		if err != nil {
			hal.Logger.Errorf("couldn't fetch karma details: %v", err)
			err = res.Reply("Sorry, I couldn't fetch karma ratings")
			return err
		}
		json.Unmarshal(val, &dev)

		devs = append(devs, dev)
	}
	hal.Logger.Debugf("Value of devs: %#v", devs)

	if len(devs) == 0 {
	}

	sort.Sort(sort.Reverse(ByKarma(devs)))
	hal.Logger.Debugf("Sorted devs list: %#v", devs)

	ranking := []string{
		"Current rankings:\n",
	}

	for _, val := range devs {
		ranking = append(ranking, val.String())
	}

	var msg string
	for _, m := range ranking {
		msg = msg + m
	}
	msg, err := sprungeSend(msg)
	if err != nil {
		hal.Logger.Error(err)
		err = res.Reply("Sorry, I had some unexpected difficulties posting the rankings")
		return err
	}

	err = res.Reply(fmt.Sprintf("The current rankings can be found here %s", msg))
	return err
})

func updateDevsWithKarma(name string, r *hal.Response) {
	names := getDevKarmaList(r)
	if !inSlice(names, name) {
		if len(names) == 0 || names[0] == "" {
			names[0] = name
		} else {
			names = append(names, name)
		}
	}
	list := strings.Join(names, ",")
	r.Robot.Store.Set("KarmaList", []byte(list))
}

func getDevKarmaList(r *hal.Response) []string {
	val, _ := r.Robot.Store.Get("KarmaList")
	list := string(val)
	names := strings.Split(list, ",")

	return names
}

func inSlice(n []string, name string) bool {
	for _, val := range n {
		if val == name {
			return true
		}
	}
	return false
}
