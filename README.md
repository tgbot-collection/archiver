# archiver

[![codecov](https://codecov.io/gh/tgbot-collection/archiver/branch/master/graph/badge.svg?token=ELEIXK6QAR)](https://codecov.io/gh/tgbot-collection/archiver)
[![docker image builder](https://github.com/tgbot-collection/archiver/actions/workflows/builder.yaml/badge.svg)](https://github.com/tgbot-collection/archiver/actions/workflows/builder.yaml)
[![Docker Pulls](https://img.shields.io/docker/pulls/bennythink/archiver)](https://hub.docker.com/r/bennythink/archiver)

🗜 Internet Archive: Wayback Machine bot

[Wayback Machine](https://t.me/wayback_machine_bot)

# WayBack Machine bot

Send any link to this bot, and it will save it to [Internet Archive](https://archive.org/).

A screenshot of this page will be sent to you as well.

# screenshots

![](assets/1.png)

![](assets/2.png)

# commands

```
start - start Wayback Machine bot
about - about this bot
ping - ping server
```

# Features

* create snapshot

# Privacy notice

This bot **WILL NEVER** collect your user id, username, last name, first name, url or anything that could be used to
track you.

This bot won't save any personal information, neither in database nor in log.

Anything that you sent to this bot is confidential from the bot's side - even your url is omitted from log system.

I value your privacy, and I know it's difficult to fight against surveillance, injustice and censorship.

> Remember, remember the Fifth of November,
> The Gunpowder Treason and Plot,
> I know of no reason
> Why the Gunpowder Treason
> Should ever be forgot.

# TODO

- [x] show snapshot result

# Usage

## Build project locally

```bash
git clone https://github.com/tgbot-collection/archiver
cd archiver
go build .
TOKEN=13245 DRIVER=$(which chromedriver) ./archiver
```

## docker

```bash
docker run -e DRIVER=$(which chromedriver) TOKEN=1234 bennythink/archiver
```

## prebuilt binary
Download it from release, and then
```shell
TOKEN=1234 DRIVER=$(which chromedriver) ./archiver
```

# License

Apache License
Version 2.0