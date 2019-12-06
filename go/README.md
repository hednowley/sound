# sound

[![Build Status](https://travis-ci.org/hednowley/sound.svg?branch=master)](https://travis-ci.org/hednowley/sound)
[![Go Report Card](https://goreportcard.com/badge/github.com/hednowley/sound)](https://goreportcard.com/report/github.com/hednowley/sound)

A snappy music server written in Go using the [Subsonic](http://www.subsonic.org) API.

-   Scan file systems and [beets](https://beets.io) databases for your music
-   Browse and stream music
-   Manage playlists
-   View cover art

_sound_ targets version 1.16.1 of the Subsonic API.

## Prerequisites

You'll need

-   Go 1.11.x
-   A [PostgreSQL](https://www.postgresql.org/) database

## Building

```shell
$ go get
$ go build
```

[Cross compiling](https://golangcookbook.com/chapters/running/cross-compiling/) is possible if you want to run _sound_ on a different device to the one you're building on.

## Testing

Before running tests, create an empty database called `sound_test` owned by a user `sound` with password `sound`.

```shell
$ go test -p 1 ./...
```

## Running

-   Make a Postgres database and user for the application use.
-   Copy your built binary along with [config.yaml](config.yaml) to your deployment directory. Configure _sound_ by editing `config.yaml`, following the comments inside the file.
-   Run _sound_ as a standalone web service or behind a proxy such as [Nginx](https://www.nginx.com/).

## Connecting with a Subsonic client

You can connect a Subsonic client to _sound_ by pointing the client to the `/subsonic` subpath (so if _sound_ is hosted at `http://www.mywebsite.com/sound` you would point your client to `http://www.mywebsite.com/sound/subsonic`).

## Roadmap

-   Integrate the Elm [front-end ](https://github.com/hednowley/sound-ui-elm)
-   Transcoding
-   FastCGI
-   Scheduling
-   Starring
-   Podcasts
-   User roles
-   More music providers
