package cfg

import (
	"time"

	rt "github.com/papaburgs/almagest/pkg/redistools"
)

type ParamType uint

const (
	LogLevel ParamType = iota
	// OtherP is just another parameter used for placeholder
	OtherP
)

func SetParam(arc *rt.AlmagestRedisClient, p ParamType, value string) error {
	k, d := getKeyDefault(p)
	if value == "" {
		value = d
	}
	return arc.Set(k, value, time.Duration(0))
}

func GetParam(arc *rt.AlmagestRedisClient, p ParamType) string {
	k, d := getKeyDefault(p)
	value, err := arc.Get(k)
	if err != nil {
		return d
	}
	return value
}

func getKeyDefault(param ParamType) (key string, def string) {
	switch param {
	case LogLevel:
		key = "almagest|config|loglevel"
		def = "info"

	case OtherP:
		key = "almagest|config|otherp"
		def = ""
	}
	return
}
