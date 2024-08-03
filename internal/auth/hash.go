package auth

import (
	"github.com/matthewhartstonge/argon2"
	"github.com/qxuken/short/internal/config"
)

func Verify(conf *config.Config, token []byte) (bool, error) {
	return argon2.VerifyEncoded(token, conf.AdminToken)
}

func GeneratePHCHash(token []byte) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(token))
	return string(encoded), err
}
