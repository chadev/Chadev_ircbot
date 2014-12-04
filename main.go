package main

import (
	"os"

	"github.com/danryan/hal"
	// _ "github.com/danryan/hal/adapter/irc"
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
