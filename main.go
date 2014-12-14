package main

import (
	"os"

	"github.com/danryan/hal"
	_ "github.com/danryan/hal/adapter/irc"
	_ "github.com/danryan/hal/adapter/shell"
	_ "github.com/danryan/hal/store/memory"
)

var pingHandler = hal.Hear(`.ping`, func(res *hal.Response) error {
	return res.Send("PONG")
})

var fooHandler = hal.Hear(`.foo`, func(res *hal.Response) error {
	return res.Send("BAR")
})

var synHandler = hal.Hear(`.SYN`, func(res *hal.Response) error {
	return res.Send("ACK")
})

var helpHandler = hal.Hear(`.help`, func(res *hal.Response) error {
	helpMsg := `HAL Chadev IRC Edition
Supported commands:
.ping      - Causes HAL to reply with a PONG
.foo       - Causes HAL to reply with a BAR
.SYN       - Causes HAL to reply with an ACK
.tableflip - ...
.events    - Get next 7 events from the Chadev calendar
.help      - Displays this message`

	return res.Send(helpMsg)
})

func run() int {
	robot, err := hal.NewRobot()
	if err != nil {
		hal.Logger.Error(err)
		return 1
	}

	robot.Handle(
		pingHandler,
		fooHandler,
		tableFlipHandler,
		eventHandler,
		synHandler,
		helpHandler,
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
