package internal

import (
	"log"
	"net"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug        bool    `default:"false"`
	DataPath     string  `envconfig:"DATA_PATH" default:"./tmp"`
	Port         int     `default:"8080"`
	Bind         net.IP  `default:"127.0.0.1"`
	PublicUrl    url.URL `envconfig:"PUBLIC_URL" required:"true"`
	PublicUrlStr string  `ignore:"true"`
	HandleLen    int     `envconfig:"HANDLE_LEN" default:"5"`
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
