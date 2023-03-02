#!/usr/bin/env bash

set -e
ctr=$(buildah from fedora)
mnt=$(buildah mount $ctr)

cp build/discord/discord-bot $mnt/discord-bot
chmod +x $mnt/discord-bot

buildah config --entrypoint /discord-bot $ctr
buildah umount $ctr
buildah commit $ctr discord-bot
