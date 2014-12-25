// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/danryan/hal"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type Chadevs struct {
	Devs []Dev `json:"devs"`
}

type Dev struct {
	Github     string `json:"github"`
	GravatarId string `json:"gravatar-id"`
	Name       string `json:"name"`
	Urls       []Url  `json"urls"`
}

type Url struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func getChadevs() (Chadevs, error) {
	url := "http://chadev.github.io/devs.json"
	res, err := http.Get(url)
	if err != nil {
		return Chadevs{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Chadevs{}, err
	}
	var data Chadevs
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Chadevs{}, err
	}
	return data, nil
}

func devsCount(d Chadevs) int {
	return len(d.Devs)
}

func devsListAll(d Chadevs) []string {
	cnt := devsCount(d)
	lst := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		lst[i] = d.Devs[i].Name
	}
	return lst
}

func findDev(d Chadevs, name string) (Dev, bool) {
	cleanedname := strings.ToUpper(strings.TrimSpace(name))
	distance := int(^uint(0) >> 1) // initialize to "infinity"
	idx := 0
	for i := 0; i < devsCount(d); i++ {
		cleaneddev := strings.ToUpper(d.Devs[i].Name)
		if cleanedname == cleaneddev {
			return d.Devs[i], true
		}
		newdistance := levenshtein.DistanceForStrings([]rune(cleanedname), []rune(cleaneddev), levenshtein.DefaultOptions)
		if newdistance < distance {
			distance = newdistance
			idx = i
		}
	}

	return d.Devs[idx], false
}

func devUrlsMessage(d Dev) string {
	cnt := len(d.Urls)
	lst := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		lst[i] = fmt.Sprintf("They have a %s. The URL is: %s.", d.Urls[i].Name, d.Urls[i].Url)
	}

	return strings.Join(lst, " ")
}

func devGravatarUrl(d Dev, size int) string {
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s.jpg?s=%v", d.GravatarId, size)
}

func devGravatarMessage(d Dev, size int) string {
	if d.GravatarId == "" {
		return "Oh no! We don't have their gravatar. What they may or may not look like is a total mystery!"
	} else {
		return fmt.Sprintf("This is their gravatar, which may or may not look like them: %s.", devGravatarUrl(d, size))
	}
}

func devGithubMessage(d Dev) string {
	if d.Github == "" {
		return "Oh no! We don't know who they are on GitHub, so we can't see where their code go!"
	} else {
		return fmt.Sprintf("Their Github is https://github.com/%s.", d.Github)
	}
}

var chadevCountHandler = hear(`chadevs count`, func(res *hal.Response) error {
	chadevs, err := getChadevs()
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("Unable to get chadevs count: %v\n", err))
		return res.Send("Can't get count. Try again maybe?")
	}

	return res.Send(fmt.Sprintf("There are currently %s chadevs.", strconv.Itoa(devsCount(chadevs))))
})

var chadevListAllHandler = hear(`chadevs all`, func(res *hal.Response) error {
	chadevs, err := getChadevs()
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("Unable to get list of chadevs: %v\n", err))
		return res.Send("Can't get list of chadevs. Try again maybe?")
	}

	lst := strings.Join(devsListAll(chadevs), ", ")

	return res.Send(fmt.Sprintf("The following people are members of Chadev: %s", lst))
})

var chadevInfoHandler = hear(`chadevs info (.+)`, func(res *hal.Response) error {
	chadevs, err := getChadevs()
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("Unable to get info of chadever: %v\n", err))
		return res.Send("Can't get info of chadev member. Try again maybe?")
	}

	dev, exactmatch := findDev(chadevs, res.Match[1])
	if !exactmatch {
		return res.Send(fmt.Sprintf("Didn't find %s. Did you mean %s?", res.Match[1], dev.Name))
	} else {
		return res.Send(fmt.Sprintf("Ah, %s! %s %s %s", dev.Name, devGravatarMessage(dev, 200), devGithubMessage(dev), devUrlsMessage(dev)))
	}
})
