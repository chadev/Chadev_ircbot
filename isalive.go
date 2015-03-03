// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/danryan/hal"
)

var isAliveHandler = hear(`is ([A-Za-z0-9\-_\{\}\[\]\\\s]+) alive`, "is (name) alive", "Find out if a user is alive", func(res *hal.Response) error {
	name := res.Match[1]
	res.Send(fmt.Sprintf("I can't find %s's heartbeat..", name))
	res.Send("But let's not jump to conclusions")
	return nil
})
