package shortener

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"strings"

	"github.com/qxuken/short/internal/db"
)

const (
	RETRIES int = 4
)

var (
	ALPHABET           []byte  = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-")
	ALPHABET_FULL_LEN  big.Int = *big.NewInt(int64(len(ALPHABET)))
	ALPHABET_SHORT_LEN big.Int = *big.NewInt(int64(len(ALPHABET) - 2))
)

func ShortUrlWithLen(length int) string {
	var res strings.Builder

	for i := 1; i < length; i++ {
		alphabetIndex, _ := rand.Int(rand.Reader, &ALPHABET_FULL_LEN)
		res.WriteByte(ALPHABET[alphabetIndex.Int64()])
	}
	alphabetIndex, _ := rand.Int(rand.Reader, &ALPHABET_SHORT_LEN)
	res.WriteByte(ALPHABET[alphabetIndex.Int64()])

	return res.String()

}

func ShortUrlChecked(db db.DB, handleLen int) (string, error) {
	for range RETRIES {
		short := ShortUrlWithLen(handleLen)
		_, err := db.GetLink(short)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return short, nil
			}
			return "", nil
		}
	}
	return "", errors.New("Cannot create unique handle")
}
