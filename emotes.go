package main

import "github.com/danryan/hal"

var tableFlipHandler = hal.Hear(`.tableflip`, func(res *hal.Response) error {
	return res.Send(`(╯°□°）╯︵ ┻━┻`)
})
