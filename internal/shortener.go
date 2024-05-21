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

func shortUrl(url []byte, l int) string {
	hasher := fnv.New64a()
	hasher.Write(url)

	hash := hasher.Sum(nil)[:l]

	return strings.ToLower(base32.StdEncoding.EncodeToString(hash))

}

func ShortUrl(url string) string {
	return shortUrl([]byte(url), 5)
}

func ShortUrlChecked(db db.DB, url string) (string, error) {
	buf := []byte(url)
	for l := range 5 {
		for range 5 {
			short := shortUrl(buf, 5+l)
			_, err := db.GetLink(short)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return short, nil
				}
				return "", nil
			}
			rand.Read(buf)
		}
	}
	return "", errors.New("Cannot create unique handle for " + url)
}
