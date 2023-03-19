GO_BUILD = go build

all: discord-bot api watchdog

generic-build: cmd/$(BOT)/gitc.txt
	mkdir -p build/$(BOT)
	$(GO_BUILD) -o build/$(BOT)/$(BOT) cmd/$(BOT)/*.go

cmd/$(BOT)/gitc.txt: gitc
	ln --force gitc.txt cmd/$(BOT)/gitc.txt

image: generic-build 
	podman unshare init/build-$(BOT).sh 

discord-bot:
	BOT=discord-bot make image

api:
	BOT=api make image

watchdog:
	BOT=watchdog make image

gitc:
	git rev-parse --short HEAD >gitc.txt

clean:
	rm -rf build
	find . -name gitc.txt -delete
	buildah rm --all

PHONY: all clean generic-build build-discord
