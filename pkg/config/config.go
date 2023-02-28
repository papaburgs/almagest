package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-redis/redis"
)

var (
	Config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
	BotID     string
	Channels  map[string]string
	RedisDB   redis.Client
}

func ReadConfig(filename string) error {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Config)
	return err

}
