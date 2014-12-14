// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/danryan/hal"

var tableFlipHandler = hal.Hear(`.tableflip`, func(res *hal.Response) error {
	return res.Send(`(╯°□°）╯︵ ┻━┻`)
})
