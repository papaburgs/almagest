package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	Config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	BotID     string `json:"BotID"`
}

func ReadConfig(filename string) error {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &Config)

	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil

}
