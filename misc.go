// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/danryan/hal"
)

var fatherHandler = hear(`your father was (a|an) (.*)`, "", "", func(res *hal.Response) error {
	return res.Send(fmt.Sprintf("%s, well your father was a hampster, and your mother spelled of elderberries!", res.UserName()))
})
