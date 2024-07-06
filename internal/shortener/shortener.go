package shortener

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"errors"
	"strings"

	"github.com/qxuken/short/internal/db"
)

const (
	RETRIES     int = 5
	DEFAULT_LEN int = 5
	MAX_LEN     int = 15
)

func shortUrl(l int) string {
	buf := make([]byte, l)
	rand.Read(buf)

	return strings.ToLower(base32.StdEncoding.EncodeToString(buf))

}

func ShortUrl() string {
	return shortUrl(5)
}

func ShortUrlChecked(db db.DB) (string, error) {
	for l := DEFAULT_LEN; l <= MAX_LEN; l++ {
		for r := 0; r < RETRIES; r++ {
			short := shortUrl(l)
			_, err := db.GetLink(short)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return short, nil
				}
				return "", nil
			}
		}
	}
	return "", errors.New("Cannot create unique handle")
}
