package internal

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
	"hash/fnv"
	"strings"

	"github.com/qxuken/short/internal/db"
)

func shortUrl(url []byte) string {
	hasher := fnv.New64a()
	hasher.Write(url)

	hash := hasher.Sum(nil)[:5]

	return strings.ToLower(base32.StdEncoding.EncodeToString(hash))

}

func ShortUrl(url string) string {
	return shortUrl([]byte(url))
}

func ShortUrlChecked(db db.DB, url string) (string, error) {
	buf := []byte(url)
	for range 20 {
		short := shortUrl(buf)
		_, err := db.GetLink(short)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return short, nil
			}
			return "", nil
		}
		rand.Read(buf)
	}
	return "", errors.New("Cannot create unique handle for " + url)
}
