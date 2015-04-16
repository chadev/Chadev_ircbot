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

// Groups contains data on the various dev groups.
type Groups struct {
	// Name of the group
	Name string `json:"name"`
	// URL for the group page/meetup page
	URL string `json:"url"`
}

var groupListHandler = hear(`(groups|meetups) list`, "(groups|meetups) list", "Lists all groups that are known to Ash", func(res *hal.Response) error {
	groups, err := res.Robot.Store.Get("GROUPS")
	if err != nil {
		res.Send("I am currently unaware of any groups, try adding some")
		return err
	}

	var g []Groups
	err = json.Unmarshal(groups, &g)
	if err != nil {
		hal.Logger.Errorf("error parsing JSON: %v", err)
		return res.Send("I had an error parsing the groups")
	}

	var gn []string
	for _, val := range g {
		gn = append(gn, val.Name)
	}
	names := strings.Join(gn, " ")

	return res.Send(fmt.Sprintf("Here is a list of groups: %s", names))
})

var groupAddHandler = hear(`(groups|meetups) add (.+): (.+)`, "(groups|meetups) add [group name]: [group url]", "Adds a new group to Ash", func(res *hal.Response) error {
	name := res.Match[2]
	url := res.Match[3]

	if name == "" {
		hal.Logger.Warn("no group name given")
		return res.Send("I need a name for the group to add it.")
	}
	if url == "" {
		hal.Logger.Warn("no group url given")
		return res.Send("I need the url for the groups webpage or meetup group")
	}

	var g []Groups
	groups, err := res.Robot.Store.Get("GROUPS")
	if len(groups) > 0 {
		err := json.Unmarshal(groups, &g)
		if err != nil {
			hal.Logger.Errorf("faild to parse json: %v", err)
			return res.Send("Failed to parse groups list")
		}
	}

	g = append(g, Groups{Name: name, URL: url})
	groups, err = json.Marshal(g)
	if err != nil {
		hal.Logger.Errorf("faild to build json: %v", err)
		return res.Send("Failed write updated groups list")
	}
	err = res.Robot.Store.Set("GROUPS", groups)
	if err != nil {
		hal.Logger.Error(err)
		return res.Send("Failed writing to the datastore")
	}

	return res.Send("Added new group")
})
