package main

import (
	"fmt"
	"log"

	"github.com/papaburgs/almagest/pkg/config"
)

func main() {
	log.Print("Reading config file from ./config.json")
	err := config.ReadConfig("config.json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	// hold application up indefinatly
	<-make(chan struct{})
	return
}
