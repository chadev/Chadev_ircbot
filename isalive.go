// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/danryan/hal"
)

var isAliveHandler = hear(`is ([A-Za-z0-9\-_\{\}\[\]\\]+) alive`, "is (name) alive", "Find out if a user is alive", func(res *hal.Response) error {
	name := res.Match[1]
	return res.Send(fmt.Sprintf("I can't find %s's heartbeat.. \nBut let's not jump to conclusions", name))
})
