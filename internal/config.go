package internal

import (
	"log"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug        bool    `default:"false"`
	Port         int     `default:"8080"`
	PublicUrl    url.URL `envconfig:"PUBLIC_URL" required:"true"`
	PublicUrlStr string  `ignore:"true"`
}

func LoadConfig() *Config {
	var s Config
	err := envconfig.Process("shortiq", &s)
	s.PublicUrlStr = s.PublicUrl.String()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}
