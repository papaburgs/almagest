all: discord-bot

discord-bot: bin/discord-bot discord-gitc
	podman unshare init/build-discord-bot.sh

bin/discord-bot:
	mkdir -p build/discord
	go build -o build/discord/discord-bot cmd/discord-bot/*.go

clean:
	rm -rf build

cmd/discord-bot/gitc.txt: gitc
	ln gitc.txt cmd/discord-bot/gitc.txt

cmd/api/gitc.txt: gitc
	ln gitc.txt cmd/api/gitc.txt

gitc:
	git rev-parse --short HEAD >gitc.txt


PHONY: all clean discord-bot gitc
