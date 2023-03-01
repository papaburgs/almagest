
all: discord-bot

discord-bot: bin/discord-bot
	podman unshare init/build-discord-bot.sh

bin/discord-bot:
	go build -o bin/discord-bot cmd/discord-bot/*.go

clean:
	rm -f bin/discord-bot


PHONY: all clean discord-bot
