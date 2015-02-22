# Chadev IRC bot
[![Build Status](https://travis-ci.org/chadev/Chadev_ircbot.svg)](https://travis-ci.org/chadev/Chadev_ircbot)

This is a custom IRC bot built for the #chadev channel.  It is based on the [HAL](https://github.com/danryan/hal) framework.

## Installation

The bot is written in Go, and reqiures [Go 1 and higher](http://golang.org/doc/install).  To build run the following:

    go get github.com/chadev/Chadev_ircbot

Currently this is unbuildable with Go 1.4, due to the way C and CGO are handled.  See the [Go 1.4 release notes](http://golang.org/doc/go1.4#swig) and [Issue #39 on HAL's issue tracker](https://github.com/danryan/hal/issues/39) for more details.

### Running tests

With the source code downloaded, the unit tests can be ran at anytime with the following:

    go test .

For more details the tests can be can verbosely with:

    go test -v .

The tests can be can at the sametime as installing the bot with the ```-t``` flag

    go get -t github.com/chadev/Chadev_ircbot

### Redis and persistent storage

By default the bot will use a memory storage.  This is fine for testing or development.
However, for persistent storage we use Redis.  Redis is available through most
package managers (apt, brew, etc.).  Once it's installed you can simply run the
redis server like so:

    $ redis-server

You will also need to supply the proper environment variables (found below)

## Configuring

The bot is configured purely by system evironmental variables.  To work with the Google Calendar API, this also requires Oauth2 credentuals from Google. To set those up follow the directions [found here](https://developers.google.com/accounts/docs/OAuth2ForDevices).

### Required environtment variables

ENV Variable | Values
-------------|-------
HAL\_ADAPTER | "shell" or "irc"
HAL\_IRC\_USER | username (string)
HAL\_IRC\_NICK | nickname (string)
HAL\_IRC\_SERVER | URL (string)
HAL\_IRC\_CHANNELS | comma seperated list of channels
HAL\_STORE | "redis" or "memory" (defaults to memory)
HAL\_REDIS\_URL | host:port (defaults to localhost:6379)
CHADEV\_TOKEN | Google Oauth2 refresh token (string)
CHADEV\_ID | Google Oauth2 Client ID (string)
CHADEV\_SECRET | Google Oauth2 Client Secret (string)
CHADEV\_MEETUP | Meetup API Key (string)
CHADEV\_PASTEBIN | Pastebin API Key (string)

## Running the bot

Running the bot is simple

    $ chadev_ircbot

This will start up the bot using whatever the environment variables are set to.  When testing, setting the adapter HAL uses to `"shell"` is helpful. This can be done at launch like so:

    $ export HAL_ADAPTER="shell"; chadev_ircbot

## Usage

All commands use the "noun verb" syntax, the noun is the name of the bot (currently "Ash").  An example of this would be ```Ash ping``` to send a ping to the bot.

Command | Details
--------|---------
events | Gets next seven events from the Chadev calendar
foo    | Causes Ash to reply with a BAR
fb n   | Return the result of FizzBuzz for n
help   | Displays the help message
issue  | Returns the URL for the issue queue for the given CHadev project
ping   | Causes Ash to reply with PONG
recall `key` | Causes the bot to read back a stored note
remember `key`: `note` | Tells the to remember something
source | Returns the URL for the given Chadev project
SYN    | Causes Ash to reply with ACK
tableflip | Flips some table
cageme | Sends Nic Cage to infiltrate your brain
who is `username` | Tells you who a user is
`username` is `description` | Tells Ash who that user is
chadevs count | Count of all members of Chadev
chadevs all | List all members of Chadev
chadevs info `full name` | Get info about Chadev member or will try to guess the name you meant
what \[are\|is\] (query) | Has the bot search for something, return the URL for the results
give me some music | Return list of music playlists popular in the community
is today dev lunch day? | Returns if today is Dev Lunch day, if so gives details on it
tell me about the next talk | Returns details for the next listed Chadev Lunch talk
devlunch url (date) (url) | Set live stream url for the dev lunch talks
link to devlunch | Returns the link to the dev lunch live stream

## License

Chadev IRC bot is licensed under the BSD 3-clause license.
