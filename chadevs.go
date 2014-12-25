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
)

type Chadevs struct {
	Devs []struct {
		Github     string `json:"github"`
		GravatarId string `json:"gravatar-id"`
		Name       string `json:"name"`
		Urls       []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"urls"`
	} `json:"devs"`
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

var chadevCountHandler = hear(`chadevs count`, func(res *hal.Response) error {
	chadevs, err := getChadevs()
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("Unable to get chadevs count: %v\n", err))
		return res.Send("Can't get count. Try again maybe?")
	}

	return res.Send("There are currently " + strconv.Itoa(devsCount(chadevs)) + " chadevs.")
})

var chadevListAllHandler = hear(`chadevs all`, func(res *hal.Response) error {
	chadevs, err := getChadevs()
	if err != nil {
		hal.Logger.Error(fmt.Sprintf("Unable to get list of chadevs: %v\n", err))
		return res.Send("Can't get list of chadevs. Try again maybe?")
	}

	lst := strings.Join(devsListAll(chadevs), ", ")

	return res.Send("The following people are members of Chadev: " + lst)
})
