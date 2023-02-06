# Discord chat bot
Eventually, be able to pull messages off of redis and post to discord.

## Icon
💬

### Key Namespaces
#### Message Posting
Message can be posted here and they will be sent to discord.
timestamp is there for sorting and hash for uniquness
```
almagest:discord:post:msg-<timestamp>-<hash>
```

#### Control Messages
Not sure what this will be used for but might be used
```
almagest:discord:control:msg-<timestamp>-<hash>
```
