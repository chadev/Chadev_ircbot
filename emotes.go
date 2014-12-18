// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"

	"github.com/danryan/hal"
)

var tableFlipHandler = hal.Hear(listenName+` tableflip`, func(res *hal.Response) error {
	num := rand.Int()
	switch {
	case num%15 == 0:
		return res.Send(`the table flipped you! ノ┬─┬ノ ︵ ( \o°o)\`)
	case num%3 == 0:
		return res.Send("(ノಠ益ಠ)ノ彡┻━┻")
	case num%5 == 0:
		return res.Send("you set the table down ┬─┬ノ( º _ ºノ)")
	default:
		return res.Send(`(╯°□°）╯︵ ┻━┻`)
	}
})
