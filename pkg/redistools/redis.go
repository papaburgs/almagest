package redistools

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

type messageNSType int8

const (
	Post messageNSType = iota
	Control
)

var messageNSMap = map[messageNSType]string{
	Post:    "post",
	Control: "control",
}

func New() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func SortedKeys(rdb *redis.Client, prefix string) []string {
	ctx := context.Background()
	keysList, err := rdb.Keys(ctx, prefix+"*").Result()
	if err != nil {
		log.Printf("Error getting keys: %s", err)
		return []string{}
	}

	sort.Strings(keysList)
	return keysList
}

// makeKey appends timestamp and hash
// content should not end in the delimeter
func appendToKey(content string) string {
	ts := time.Now().Unix()
	id := uuid.New()

	return fmt.Sprintf("%s:%v:%s", content, ts, id)
}

// DiscordKey returns a namespaced and hashed key
// set second argument to true to just get the prefix
func DiscordKey(ns messageNSType, justPrefix bool) string {
	var content string
	if ns == Post {
		content = "almagest:discord:post:msg"
	}
	if ns == Control {
		content = "almagest:discord:control:msg"
	}
	if justPrefix {
		return content
	}
	return appendToKey(content)
}
