# Chadev IRC bot
[![Build Status](https://travis-ci.org/chadev/Chadev_ircbot.svg)](https://travis-ci.org/chadev/Chadev_ircbot)

This is a custom IRC bot built for the #chadev channel.  It is based on the [HAL](https://github.com/danryan/hal) framework.

## Installation

The bot is written in Go, and reqiures [Go 1 and higher](http://golang.org/doc/install).  To build run the following:

    go get github.com/chadev/Chadev_ircbot

Currently this is unbuildable with Go 1.4, due to the way C and CGO are handled.  See the [Go 1.4 release notes](http://golang.org/doc/go1.4#swig) and [Issue #39 on HAL's issue tracker](https://github.com/danryan/hal/issues/39) for more details.

## Configuring

The bot is configured purly by system evironmental variables.  To work with the Google Calendar API, this also requires Oauth2 credentuals from Google.
To set those up follow the directions [found here](https://developers.google.com/accounts/docs/OAuth2ForDevices).

### Required environtment variables

ENV Variable | Values
-------------|-------
HAL\_ADAPTER | "shell" or "irc"
HAL\_IRC\_USER | username (string)
HAL\_IRC\_NICK | nickname (string)
HAL\_IRC\_SERVER | URL (string)
HAL\_IRC\_CHANNELS | comma seperated list of channels
CHADEV\_TOKEN | Google Oauth2 refresh token (string)
CHADEV\_ID | Google Oauth2 Client ID (string)
CHADEV\_SECRET | Google Oauth2 Client Secret (string)

## Running the bot

Running the bot is simple

    $ chadev_ircbot

This will start up the bot using what ever the environment variables are set to.  When testing switching the adapter HAL will use to shell is helpful,
this can be done at launch like so:

    $ export HAL_ADAPTER="shell"; chadev_ircbot

## Usage

Currently the help output is broken, see [Issue #1](https://github.com/chadev/Chadev_ircbot/issues/1) for more details.  Until that is worked out here is a list of supported commands.

All commands use the "noun verb" syntax, the noun is the name of the bot (currently "Ash").  An example of this would be ```Ash ping``` to send a ping to the bot.

Command | Details
--------|---------
.events | Get next seven events from the Chadev calendar
.foo    | Causes HAL to reply with a BAR
.help   | Displays the help message
.ping   | Causes HAL to reply with PONG
.SYN    | causes HAL to reply with ACK
.tableflip | Flips a table

## License

Chadev IRC bot is licensed under the BSD 3-clause license.
