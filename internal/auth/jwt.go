package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qxuken/short/internal/config"
)

func IssueJWT(conf *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(conf.AppSecret)
}

func VerifyJWT(conf *config.Config, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return conf.AppSecret, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
