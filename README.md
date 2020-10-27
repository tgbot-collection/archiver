# archiver
🗜 Internet Archive: Wayback Machine bot

[Wayback Machine](https://t.me/wayback_machine_bot)

# WayBack Machine bot
Send any link to this bot, and it will save it to [Internet Archive](https://archive.org/).

# commands
```
start - start Wayback Machine bot
about - about this bot
```

# Features
* create snapshot

# TODO
- [ ] show snapshot result
- [ ] history records show and clear

# Build
## General approach
```bash
git clone https://github.com/tgbot-collection/archiver
cd archiver
go build .
token=13245 ./archiver
```
## docker
```bash
docker run -e token=1234 bennythink/archiver
```
# License