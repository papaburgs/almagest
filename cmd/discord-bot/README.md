# Discord chat bot
Eventually, be able to pull messages off of redis and post to discord.

## Icon
💬

## PubSub implementation

Discord bot subscribes alamgest channel and watches for messages on its service.
It will post those messages to discod

It also watches the 32ohsix discord server and checks posted messages for items.
If it is something that matches a rule, it can action that by publishing messages to redis.


## Install
* `make discord-bot` will make an image in localhost registry
* can run it with something like: 

```
podman run --rm --secret source=DISCORD_BOT_TOKEN,type=env --pod almagest discord-bot
```

using a podman secret, where the secret is the discord token
