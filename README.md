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

<img width="1217" alt="Screenshot 2022-06-16 at 17 38 42" src="https://user-images.githubusercontent.com/2153895/174094944-376cd186-03b8-4cf9-8828-4d5256676e7c.png">

## Usage

Find bot https://t.me/ShoppingCatBot and add him to your chats.

To run your own server, just pull docker image, pass your Telegram token and run container:

```shell
docker pull ghcr.io/egregors/shoppingcatbot/scbot:latest
docker run                                      \
  --rm                                          \
  -d                                            \
  -v /tmp/dumps:/dumps                          \
  -e SCBOT_TG_TOKEN="YOUR-TELEGRAM-BOT-TOKEN"   \
  ghcr.io/egregors/shoppingcatbot/scbot
```

## Development

Use `make` to run developer's commands

```shell
Usage: make [task]

task                 help
------               ----
run                  Run local dev version
drun                 Run in docker
lint                 Lint all the stuff
test                 Run tests
docker               Build docker image
                     
help                 Show help message

```

## Contributing

Bug reports, bug fixes and new features are always welcome.
Please open issues and submit pull requests for any new code.
