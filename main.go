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

const VERSION = "2014-12-28"

// handler is an interface for objects to implement in order to respond to messages.
type handler interface {
	Handle(res *hal.Response) error
}

var helpMessages = make(map[string]string)

var pingHandler = hear(`ping`, "ping", "Causes Ash to reply with PONG", func(res *hal.Response) error {
	return res.Send("PONG")
})

var fooHandler = hear(`foo`, "foo", "Causes Ash to reply with a BAR", func(res *hal.Response) error {
	return res.Send("BAR")
})

var synHandler = hear(`SYN`, "SYN", "Causes Ash to reply with ACK", func(res *hal.Response) error {
	return res.Send("ACK")
})

var selfHandler = hear(`who are you`, "self", "", func(res *hal.Response) error {
	return res.Send("I'm Ash, the friendly #chadev bot.  I can preform a variety of tasks, and I am learning new tricks all the time.  I am open source, and pull requests are welcome!")
})

var quitHandler = hear(`(.*)+/quit(.*)+`, "quit", "", func(res *hal.Response) error {
	name := res.UserName()
	return res.Send(fmt.Sprintf("No!  Bad %s!", name))
})

var helpHandler = hear(`help`, "help", "Displays this message", func(res *hal.Response) error {
	helpMsg := []string{
		"HAL Chadev IRC Edition build: " + VERSION,
		"Supported commands:",
	}

	for command, message := range helpMessages {
		if command != "" && message != "" {
			helpMsg = append(helpMsg, command+" - "+message+"\n")
		}
	}

	for _, msg := range helpMsg {
		res.Send(msg)
		time.Sleep(100 * time.Millisecond)
	}

	return nil
})

func hear(pattern string, command string, message string, fn func(res *hal.Response) error) handler {
	addHelpMessage(command, message)
	return hal.Hear("^(?i)Ash "+pattern, fn)
}

func addHelpMessage(command string, message string) {
	helpMessages[command] = message
}

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
		whoamHandler,
		isHandler,
		selfHandler,
		quitHandler,
		fizzBuzzHandler,
		noteStoreHandler,
		noteGetHandler,
		noteRemoveHandler,
		chadevCountHandler,
		chadevListAllHandler,
		chadevInfoHandler,
		fatherHandler,
		partyHandler,
		whoBackHandler,
		whatAreHandler,
	)

	if err := robot.Run(); err != nil {
		hal.Logger.Error(err)
		return 1
	}
	return 0
}
