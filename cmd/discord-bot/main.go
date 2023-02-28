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
		fmt.Println("Could not read config file: ", err.Error())
		return
	}

	done := make(chan bool, 1)
	Start(done)

	// hold application up indefinatly
	log.Print("at done")
	<-done
	log.Print("past done")
}
