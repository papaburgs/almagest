all: discord-bot

discord-bot: bin/discord-bot
	podman unshare init/build-discord-bot.sh

bin/discord-bot:
	mkdir -p build/discord
	go build -o build/discord/discord-bot cmd/discord-bot/*.go

clean:
	rm -rf build

PHONY: all clean discord-bot
