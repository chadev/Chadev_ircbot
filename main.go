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
	_ "github.com/danryan/hal/store/redis"
)

// handler is an interface for objects to implement in order to respond to messages.
type handler interface {
	Handle(res *hal.Response) error
}

var pingHandler = hear(`ping`, func(res *hal.Response) error {
	return res.Send("PONG")
})

var fooHandler = hear(`foo`, func(res *hal.Response) error {
	return res.Send("BAR")
})

var synHandler = hear(`SYN`, func(res *hal.Response) error {
	return res.Send("ACK")
})

var selfHandler = hear(`who are you`, func(res *hal.Response) error {
	return res.Send("I'm Ash, the friendly #chadev bot.  I can preform a variety of tasks, and I am learning new tricks all the time.  I am open source, and pull requests are welcome!")
})

var helpHandler = hear(`help`, func(res *hal.Response) error {
	helpMsg := `HAL Chadev IRC Edition
Supported commands:
events    - Get next 7 events from the Chadev calendar
foo       - Causes HAL to reply with a BAR
help      - Displays this message
issues    - Give the URL to the issue queue for the named GitHub repo
ping      - Causes HAL to reply with a PONG
source    - Give the URL to the named GitHub repo
SYN       - Causes HAL to reply with an ACK
tableflip - ...
cageme    - Sends Nic Cage to infiltrate your brain
who is    - Find out who a user is
(name) is - Tell HAL who the user is`

	return res.Send(helpMsg)
})

func main() {
	os.Exit(run())
}

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
		issueHandler,
		cageMeHandler,
		whoisHandler,
		isHandler,
		selfHandler,
		whoamHandler,
	)

	if err := robot.Run(); err != nil {
		hal.Logger.Error(err)
		return 1
	}
	return 0
}

func hear(pattern string, fn func(res *hal.Response) error) handler {
	return hal.Hear("^(?i)Ash "+pattern, fn)
}
