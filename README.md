# ShoppingCatBot

Shopping list telegram Bot

- 🐱 Tiny docker image
- 📝 Supports: adding, wiping and showing items
- 🧘‍♀️ Add the bot to any chat or use by yourself

> **Warning**
> Bot is in active development now, so I bet you gonna lose you date a few times, until first stable release ;)

---

<div align="center">

[![Build Status](https://github.com/egregors/ShoppingCatBot/actions/workflows/go.yml/badge.svg)](https://github.com/egregors/ShoppingCatBot/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/egregors/ShoppingCatBot)](https://goreportcard.com/report/github.com/egregors/ShoppingCatBot)

</div>

## Screenshots

[//]: # (TODO: add screenshots from an iPhone)

## Usage

Find bot https://t.me/ShoppingCatBot and add him to your chats.

To run your own server, just pull docker image, pass your Telegram token and run container:

```shell
docker pull ghcr.io/egregors/shoppingcatbot/scbot:latest
docker run -d -e SCBOT_TG_TOKEN="YOUR-TELEGRAM-BOT-TOKEN" ghcr.io/egregors/shoppingcatbot/scbot
```

## Development

Use `make` to run developer's commands

```shell
Usage: make [task]

task                 help
------               ----
run                  Run dev
lint                 Lint the files
docker               Build Docker image
                     
help                 Show help message

```

## Contributing

Bug reports, bug fixes and new features are always welcome.
Please open issues and submit pull requests for any new code.
