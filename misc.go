// Copyright 2014-2015 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/danryan/hal"
)

var fatherHandler = hear(`your father was (a|an) (.*)`, "", "", func(res *hal.Response) error {
	return res.Send(fmt.Sprintf("%s, well your father was a hampster, and your mother smelled of elderberries!", res.UserName()))
})

var whoBackHandler = hear(`who has your back(\?)?`, "", "", func(res *hal.Response) error {
	return res.Send(fmt.Sprintf("%s has my back!", res.UserName()))
})
