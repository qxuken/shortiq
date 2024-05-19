package internal

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug     bool   `default:"false"`
	Port      int    `default:"8080"`
	PublicUrl string `envconfig:"PUBLIC_URL" default:"/"`
}

func LoadConfig() *Config {
	var s Config
	err := envconfig.Process("shortiq", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}
