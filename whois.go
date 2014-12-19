// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/danryan/hal"
)

var whoisHandler = hear(`who\s?is (.+)`, func(res *hal.Response) error {
	name := res.Match[1]
	key := strings.ToUpper(name)
	val, err := res.Robot.Store.Get(key)
	if err != nil {
		res.Send(fmt.Sprintf("%v is no one to me\n", name))
		return err
	}
	return res.Send(fmt.Sprintf("%v is %v", name, string(val)))
})

var isHandler = hear(`([^who].+) is (.+)`, func(res *hal.Response) error {
	name := res.Match[1]
	key := strings.ToUpper(name)
	role := res.Match[2]

	storedRoles, err := res.Robot.Store.Get(key)
	roleToStore := role
	if len(storedRoles) > 0 {
		roleToStore = roleToStore + ", " + string(storedRoles)
	}

	err = res.Robot.Store.Set(key, []byte(roleToStore))
	if err != nil {
		res.Send(fmt.Sprint("There's something wrong"))
		return err
	}
	return res.Send(fmt.Sprintf("Got it. %v is %v\n", name, role))
})
