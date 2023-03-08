# API
A service to recieve api requests that can act on calls.

## Endpoints
### Discord message
Send a message to be posted to a discord channel
* path: /api/almagest/discord/dispatch
* method: GET
* return help with example post json content

* method: POST
* content: JSON {"channel": "botspot", "message": "post me", "server": "32ohsix"}
  * at this time, server is optional as there is only one connected.

### Status
* path: /api/almagest/status/<service>

service is optional, if not provided, all known will be returned

#### Returned structure
```
[
  {
    "service": [service name],
    "healthcheck": [up/down],
    "version": [version hash]
  }
]
```


