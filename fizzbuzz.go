// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"

	"github.com/danryan/hal"
)

var fizzBuzzHandler = hear(`fb ([0-9]+)`, func(res *hal.Response) error {
	i, err := strconv.ParseInt(res.Match[1], 10, 16)
	if err != nil {
		res.Send("What are you even?")
		return err
	}
	switch {
	case i%15 == 0:
		return res.Send("FizzBuzz")
	case i%3 == 0:
		return res.Send("Fizz")
	case i%5 == 0:
		return res.Send("Buzz")
	default:
		return res.Send(res.Match[1])
	}
})
