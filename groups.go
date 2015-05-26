// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chadev/Chadev_ircbot/meetup"
	"github.com/danryan/hal"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Groups contains data on the various dev groups.
type Groups struct {
	Group []Group `json:"groups"`
}

// Group contains data from the "groups" array in the JSON object.
type Group struct {
	Name          string `json:"name"`
	GitHub        string `json:"github"`
	CodeofConduct string `json:"code-of-conduct"`
	Urls          []URLs `json:"urls"`
	meetup        string
}

// URLs contains the name and url to the group website
type URLs struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func parseGroups() (Groups, error) {
	var g Groups
	r, err := http.Get("http://chadev.com/groups.json")
	if err != nil {
		return g, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return g, err
	}

	err = json.Unmarshal(b, &g)
	if err != nil {
		return Groups{}, err
	}

	return g, nil
}

var groupListHandler = hear(`(groups|meetups) list`, "(groups|meetups) list", "Lists all groups that are known to Ash", func(res *hal.Response) error {
	g, err := parseGroups()
	if err != nil {
		hal.Logger.Errorf("failed parsing group list: %v", err)
		res.Send("Sorry, I encountered an error while parsing the groups list")
		return err
	}

	var gn []string
	for _, val := range g.Group {
		gn = append(gn, val.Name)
	}
	names := strings.Join(gn, ", ")

	return res.Send(fmt.Sprintf("Here is a list of groups: %s", names))
})

var groupDetailsHandler = hear(`(group|meetup) details (.+)`, "(group|meetup) details [group name]", "Returns details about a group", func(res *hal.Response) error {
	name := res.Match[2]

	g, err := parseGroups()
	if err != nil {
		hal.Logger.Errorf("failed parsing group list: %v", err)
		res.Send("Sorry, I encountered an error while parsing the groups list")
		return err
	}

	group := searchGroups(g, name)
	for _, u := range group.Urls {
		if u.Name == "website" || u.Name == "meetup" {
			m := parseMeetupName(u.URL)
			if m != "" {
				group.meetup = m
			}
		}
	}

	var nextEvent string
	if group.meetup != "" {
		nextEvent, err = meetup.GetNextMeetup(group.meetup)
		if err != nil {
			hal.Logger.Errorf("failed fetching event from meetup.com: %v", err)
		}
	}

	var urls []string
	for _, u := range group.Urls {
		urls = append(urls, u.URL)
	}
	ul := strings.Join(urls, ", ")
	res.Send(fmt.Sprintf("Group name: %s URL: %s", group.Name, ul))
	if nextEvent != "" {
		res.Send(nextEvent)
	}

	return nil
})

var groupRSVPHandler = hear(`(group|meetup) rsvps (.+)`, "(group|meetup) rsvps [group name]", "Gets the RSVP count for the named group's next meeting", func(res *hal.Response) error {
	name := res.Match[2]

	g, err := parseGroups()
	if err != nil {
		hal.Logger.Errorf("failed parsing group list: %v", err)
		res.Send("Sorry, I encountered an error while parsing the groups list")
		return err
	}

	group := searchGroups(g, name)
	for _, u := range group.Urls {
		if u.Name == "website" || u.Name == "meetup" {
			m := parseMeetupName(u.URL)
			if m != "" {
				group.meetup = m
			}
		}
	}

	if group.meetup == "" {
		res.Send(fmt.Sprintf("%s is using an unsupported event system, can't fetch RSVP information", group.Name))
		return nil
	}

	rsvp, err := meetup.GetMeetupRSVP(group.meetup)
	if err != nil {
		hal.Logger.Errorf("failed fetching RSVP information: %v", err)
		res.Send("I was unable to fetch the latest RSVP informaion for this group")
		return err
	}

	if rsvp != "" {
		res.Send(fmt.Sprintf("%s RSVP breakdown: %s", group.Name, rsvp))
	} else {
		res.Send("There are either no upcoming events or no RSVP for the event yet")
	}

	return nil
})

func parseMeetupName(u string) string {
	// meetup URLs are structured as www.meetup.com/(group name)
	u = strings.TrimSuffix(u, "/") // trim trailing slash if present
	parts := strings.Split(u, "/")

	return parts[len(parts)-1]
}

func searchGroups(g Groups, n string) Group {
	distance := int(^uint(0) >> 1) // nitialize to "infinity"
	var idx int
	n = strings.ToUpper(strings.TrimSpace(n))

	for i := 0; i < len(g.Group); i++ {
		cleanGroup := strings.ToUpper(strings.TrimSpace(g.Group[i].Name))
		if n == cleanGroup {
			return g.Group[i]
		}

		newdistance := levenshtein.DistanceForStrings([]rune(n), []rune(cleanGroup), levenshtein.DefaultOptions)
		if newdistance < distance {
			distance = newdistance
			idx = i
		}
	}

	return g.Group[idx]
}
