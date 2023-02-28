# Almagest
set of utility apps for home server

## Implementation
see the cmd directory for implemented features.

- Run all of these in a pod beside a standard redis container
  - might replace with an internal database later, like on bolt or badger, but for now, get all the redis features and data structures for free

## Redis
Originally planned to do this with something like Kafka or RabbitMQ, but redis is light weight and easy to setup,
and I can get some database features (like for the DNS log parser)



## Features Implemented
None

## Feature ideas
### Discord bot
bot to post to discord, also monitor a server for posts that might need to be actioned (ie, drop a message in redis for another function)

### Api to the outside
a web listener for api calls, most of the time just transfer content to redis to be actioned by something else.

#### Endpoints
* `/api/rproxy/<queue-stub>`
* `/api/redis/metrics/<queue-stub>`
* `/api/redis/publish/<queue-stub>`

Also plan to use some sort of extra path just for ob-security

main use: adds items to a pubsub queue

### Control module
* subscribe to the control channel
* track what services are up
    * send out a ping message to the control channel and listen for responses
* make artribrary messages for testing
    * might just make a file listener and if a message shows up it publishes to the channel
* pull arbitrary keys - like for metrics
    * this might be better handled with an api server

### torrent notifier
* watch for started or completed torrents, send to discord bot to alert

### DNS log parser
Not sure of all the features yet, but a 'big brother' appthat will go through CoreDNS logs and report on things (ie, this user's phone has returned).

### Pod or volume builders
for home wiki, rebuild based on release in github

## Project layout
Since all packages will either pull or subscribe to redis, will make one central redis package that does all of that.

* setup connection
* return channels for subscriptions
* provide functions to publish

Then each other service can reference that package as well as do its own thing.
Since packages (like the localhost info service) will need to communicate with other services, it won't connect directly to them,
but instead just publish a message to redis and hope it gets picked up.

Most messages passed will be simple strings, but if more data is needed, they are passed as json strings.
