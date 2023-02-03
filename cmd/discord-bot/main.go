package main

import (
	"fmt"

	"github.com/papaburgs/almagest/pkg/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	// hold application up indefinatly
	<-make(chan struct{})
	return
}
