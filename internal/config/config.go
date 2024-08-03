package config

import (
	"log"
	"net"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Verbose      bool    `default:"false"`
	DataPath     string  `envconfig:"DATA_PATH" default:"./tmp"`
	Port         int     `default:"8080"`
	Bind         net.IP  `default:"127.0.0.1"`
	PublicUrl    url.URL `envconfig:"PUBLIC_URL" required:"true"`
	PublicUrlStr string  `ignore:"true"`
	HandleLen    int     `envconfig:"HANDLE_LEN" default:"5"`
	AdminToken   []byte  `envconfig:"ADMIN_TOKEN" required:"true"`
}

func LoadConfig() *Config {
	var s Config
	err := envconfig.Process("shortiq", &s)
	s.PublicUrlStr = s.PublicUrl.String()
	if err != nil {
		log.Fatal(err.Error())
	}

	if s.Verbose {
		log.Printf("Config: %+v\n", s)
	}
	return &s
}
