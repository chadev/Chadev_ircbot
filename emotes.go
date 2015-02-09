// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"time"

	"github.com/danryan/hal"
)

var tableFlipHandler = hear(`tableflip`, "tableflip", "...", func(res *hal.Response) error {
	rand.Seed(time.Now().UTC().UnixNano())

	e := []string{
		"(ノಠ益ಠ)ノ彡┻━┻",
		`(╯°□°）╯︵ ┻━┻`,
		`(╯°□°）╯︵ <ǝlqɐʇ>`,
		`the table flipped you! ノ┬─┬ノ ︵ ( \o°o)\`,
		"┻━┻ ︵ヽ(`Д´)ﾉ︵ ┻━┻",
	}

	return res.Send(e[rand.Intn(len(e))])
})
