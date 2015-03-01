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

The bot is configured purely by system evironmental variables.  Various services that
the bot works with requires API keys or other Oauth2 credentials, see the documentaiton for
the service in question on how to get these.

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
CHADEV\_TOKEN | [Google Oauth2 refresh token](https://developers.google.com/accounts/docs/OAuth2ForDevices) (string)
CHADEV\_ID | [Google Oauth2 Client ID](https://developers.google.com/accounts/docs/OAuth2ForDevices) (string)
CHADEV\_SECRET | [Google Oauth2 Client Secret](https://developers.google.com/accounts/docs/OAuth2ForDevices) (string)
CHADEV\_MEETUP | [Meetup API Key](https://secure.meetup.com/meetup_api/key/) (string)
CHADEV\_PASTEBIN | [Pastebin API Key](http://pastebin.com/api#1) (string)

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

## Contributing

Do you want to help make Ash better?  Looking for a project to work on to help you learn/try out Go?  Then you found the correct part of this README :thumbsup:.
The point of this section of the README is to help you get started contributing to the project.

### Getting Started

The first thing you should do is get [Go 1.0 or higher](http://golang.org/doc/install) installed, making sure to set your $GOPATH  and $GOBIN to a directory that you
have read and write access to.

In BASH and ZSH this may look like:

  export GOPATH=$HOME/go
  export GOBIN=$GOPATH/bin

In CSH and TCSH this may look like:

  setenv GOPATH=$HOME/go
  setenv GOBIN=$GOPATH/bin

The ```go get``` command uses the $GOPATH to store both source code and compiled versions of all packagse you install.
The source code gets vendored in ```$GOPATH/src```  Built packages (libraries) gets stored in ```$GOPATH/pkg``` under the same directory structure as the source code.
Compiled binaries gets placed in ```$GOPATH/bin``` or ```$GOBIN``` for short, as such the $GOBIN needs to be added to your $PATH.

After that is done it is preferable to setup your IDE/text editor, while a lot of modern IDEs and text editors have native Golang support if you prefer one that
does not then chances are there is a plugin/extension for your preffered IDE/editor if extensable.

For VIM the [vim-go](https://github.com/fatih/vim-go) plugin is prefered.

It is strongly reccomended to setup your IDE/Editor to run either ```gofmt``` or ```goimports```, installable by running ```go get golang.org/x/tools/cmd/goimports```.If a pull request is submitted and neither one has been ran on your code then the pull request my not be accepted until this is done, you will be politely
asked/reminded to do so, or your pull request may be updated so that it complies with this.  While this may seem strict to new comers, this is the Go standard style
guide, for an in depth reason behind this please read the sections on formatting to both [Effective Go](http://golang.org/doc/effective_go.html#formatting) and
the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments#gofmt).

Once all the technical things are out of the way, you are now ready to clone this repository and get started.  To get a copy of the repo, fork this repo and make your
changes there and sumbit Pull Requests.  Please don't submit pull requests directly to the ```master``` branch as this branch should directly reflect the version of
code that is currently deployed.  Instead submit your pull requsets to the ```develop``` branch.  Doing so does not mean that your pull requests wont be accepted by
any means.

Lastly don't be afraid to have fun with things!

### Indicating developer intoxication levels during a commit (optional)

When authoring a commit message feel free to add your current intoxication level, to keep things simple we have a 5 beer system inplace using GitHub's :beer: and :beers: emoji.

The way this system works is as follows:

-  1 :beer: = Not so drunk/just getting started
-  5 :beer: = [Ballmer Peak](http://xkcd.com/323/)
-  6+ :beer: = WTF????  Should you even keyboard right now?  :worried:

For heavy drinking sessions feel free to subsitute :beer: with :beers:

## License

Chadev IRC bot is licensed under the BSD 3-clause license.
