// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/danryan/hal"
	_ "github.com/danryan/hal/adapter/irc"
	_ "github.com/danryan/hal/adapter/shell"
	_ "github.com/danryan/hal/store/memory"
)

var listenName = "Ash"

var pingHandler = hal.Hear(listenName+` ping`, func(res *hal.Response) error {
	return res.Send("PONG")
})

var fooHandler = hal.Hear(listenName+` foo`, func(res *hal.Response) error {
	return res.Send("BAR")
})

var synHandler = hal.Hear(listenName+` SYN`, func(res *hal.Response) error {
	return res.Send("ACK")
})

var helpHandler = hal.Hear(listenName+` help`, func(res *hal.Response) error {
	helpMsg := `HAL Chadev IRC Edition
Supported commands:
events    - Get next 7 events from the Chadev calendar
foo       - Causes HAL to reply with a BAR
help      - Displays this message
ping      - Causes HAL to reply with a PONG
source    - Give the URL to the named GitHub repo
SYN       - Causes HAL to reply with an ACK
tableflip - ...`

	return res.Send(helpMsg)
})

func run() int {
	robot, err := hal.NewRobot()
	if err != nil {
		hal.Logger.Error(err)
		return 1
	}

	robot.Handle(
		fooHandler,
		tableFlipHandler,
		eventHandler,
		synHandler,
		helpHandler,
		pingHandler,
		sourceHandler,
	)

	if err := robot.Run(); err != nil {
		hal.Logger.Error(err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
