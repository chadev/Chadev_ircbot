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
	// Meetup is the group name from the meetup URL
	// this is used for Meetup API calls.
	Meetup string `json:"meetup_name"`
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
	names := strings.Join(gn, ", ")

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

	var meetupName string
	if strings.Contains(url, "meetup.com") {
		meetupName = parseMeetupName(url)
	}

	g = append(g, Groups{Name: name, URL: url, Meetup: meetupName})
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

var groupDetailsHandler = hear(`(group|meetup) details (.+)`, "(group|meetup) details [group name]", "Returns details about a group", func(res *hal.Response) error {
	name := res.Match[1]

	var g []Groups
	groups, err := res.Robot.Store.Get("GROUPS")
	if len(groups) > 0 {
		err := json.Unmarshal(groups, &g)
		if err != nil {
			hal.Logger.Errorf("faild to parse json: %v", err)
			return res.Send("Failed to parse groups list")
		}
	}

	group := searchGroups(g, name)
	e := Groups{}
	if group == e {
		// TODO catch group not found
		hal.Logger.Warnf("no group with the name %s found", name)
		return res.Send(fmt.Sprintf("I could not find a group with the name %s", name))
	}

	return res.Send(fmt.Sprintf("Group name: %s URL: %s", group.Name, group.URL))

})

func parseMeetupName(u string) string {
	// meetup URLs are structured as www.meetup.com/(group name)
	parts := strings.Split(u, "/")

	return parts[len(parts)-1]
}

func searchGroups(g []Groups, n string) Groups {
	var group Groups
	for _, val := range g {
		if val.Name == n {
			group = val
		}
	}

	return group
}
