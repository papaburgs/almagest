# Discord chat bot
Eventually, be able to pull messages off of redis and post to discord.

## Icon
💬

## PubSub implementation

Discord bot subscribes to a pattern so it can pick up any messages that match it.
the pattern is "almagest|discord|*".

In order to get it to pick up a message and post it, the bot must get the channel as well
a properly formed publish will then look like (from redis-cli): 
```
publish almagest|discord|post|botspot "we want elton"
```
This will post to the 'botspot' channel with the provided message

#### Control Messages
Not sure what this will be used for but might be used.

One option would be to kick off a reload of config (ie, a new token)
```
almagest|discord|control|reload "true"
```
