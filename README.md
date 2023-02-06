# Almagest
set of utility apps for home server

## Implementation
see the cmd directory for implemented features.

- Run all of these in a pod beside a standard redis container
  - might replace with an internal database later, like on bolt or badger, but for now, get all the redis features and data structures for free

## Redis?
Originally planned to do this with something like Kafka or RabbitMQ, but redis is light weight and easy to setup,
and I can get some database feaures (like for the DNS log parser)

### Features Implemented
None

### Feature ideas
#### Discord bot
bot to post to discord, also monitor a server for posts that might need to be actioned (ie, drop a message in redis for another function)

#### Api to the outside
a web listener for api calls, most of the time just transfer content to redis to be actioned by something else.

### Endpoints
`/api/rproxy/<queue-stub>`

Also plan to use some sort of extra path just for ob-security

main use: pushes (`LPUSH`) the data into a redis list. These are then extracted by another process.


#### DNS log parser
Not sure of all the features yet, but a 'big brother' appthat will go through CoreDNS logs and report on things (ie, this user's phone has returned).

#### Pod or volume builders
for home wiki, rebuild based on release in github


