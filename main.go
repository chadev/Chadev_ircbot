// Copyright 2014 Chadev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

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

var quitHandler = hear(`(.*)+/quit(.*)+`, func(res *hal.Response) error {
	name := res.UserName()
	return res.Send(fmt.Sprintf("No!  Bad %s!", name))
})

var helpHandler = hear(`help`, func(res *hal.Response) error {
	helpMsg := []string{
		"HAL Chadev IRC Edition",
		"Supported commands:",
		"events - Gets next seven events from the Chadev calendar",
		"foo - Causes HAL to reply with a BAR",
		"fb n - Return the result of FizzBuzz for n",
		"help - Displays this message",
		"issue - Returns the URL to the issue queue for the given Chadev project",
		"ping - Causes HAL to reply with a PONG",
		"source - Returns the URL the the given Chadev project",
		"SYN - Causes HAL to reply with ACK",
		"tableflip - Flip some tables",
		"cageme - Sends Nic Cage to infiltrate your brain",
		"who is (username) - Tells you who a user is",
		"(username) is (role) - Tells the bot who that user is",
	}

	for _, msg := range helpMsg {
		res.Send(msg)
		time.Sleep(100 * time.Millisecond)
	}

	return nil
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
		quitHandler,
		fizzBuzzHandler,
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
