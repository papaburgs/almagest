package main

import (
	"context"
	"log"
	"time"

	"github.com/papaburgs/almagest/pkg/redistools"
	rt "github.com/papaburgs/almagest/pkg/redistools"
)

func main() {
	ctx := context.Background()
	rdb := rt.New()

	keys := rt.SortedKeys(rdb, redistools.DiscordKey(rt.Post, true))
	log.Printf("should be no keys: %s", keys)

	err := rdb.Set(ctx, rt.DiscordKey(rt.Post, false), "value", 10*time.Minute).Err()
	if err != nil {
		panic(err)
	}
	log.Println("set is complete")

	keys = rt.SortedKeys(rdb, redistools.DiscordKey(rt.Post, true))
	log.Printf("after keys: %s", keys)

}
