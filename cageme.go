// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"math/rand"

	"github.com/danryan/hal"
)

var cageMeHandler = hal.Hear(listenName+` cageme`, func(res *hal.Response) error {
	root := "http://cageme.herokuapp.com"
	num := strconv.Itoa(int(rand.Float64()*79+1))

	return res.Send(root + "/specific/" + num + ".jpeg")
})
