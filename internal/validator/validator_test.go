package validator_test

import (
	"net/url"
	"testing"

	"github.com/qxuken/short/internal/config"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/validator"
)

const MOCK_URL_STR = "https://short.iq"

var (
	MOCK_URL, _ = url.Parse(MOCK_URL_STR)
)

func createConf() *config.Config {
	config := new(config.Config)
	config.PublicUrl = *MOCK_URL
	return config
}

func createDB() *db.SqliteDB {
	db := db.ConnectSqlite3(&config.Config{}, ":memory:")
	return db
}

func TestSimpleValidRedirectUrl(t *testing.T) {
	conf := createConf()
	url := "https://github.com/qxuken"
	res := validator.ValidateRedirectUrl(conf, url, true)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestComplexValidRedirectUrl(t *testing.T) {
	conf := createConf()
	url := "https://github.com/qxuken/key-val?query=some+hey"
	res := validator.ValidateRedirectUrl(conf, url, true)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestProtocollesInvalidRedirectUrl(t *testing.T) {
	conf := createConf()
	url := "github.com/qxuken/key-val?query=some+hey"
	res := validator.ValidateRedirectUrl(conf, url, true)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestShortInvalidRedirectUrl(t *testing.T) {
	conf := createConf()
	url := "/qxuken"
	res := validator.ValidateRedirectUrl(conf, url, true)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestMatchingInvalidRedirectUrl(t *testing.T) {
	conf := createConf()
	url := MOCK_URL_STR + "/test"
	res := validator.ValidateRedirectUrl(conf, url, true)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestSimpleValidShortUrl(t *testing.T) {
	db := createDB()
	url := "qxuken"
	res := validator.ValidateShortHandle(db, url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestComplexValidShortUrl(t *testing.T) {
	db := createDB()
	url := "AaZz09qxuke23n-sadfashfsD_1sdasfdf"
	res := validator.ValidateShortHandle(db, url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestShortInvalidShortUrl(t *testing.T) {
	db := createDB()
	url := "qxuk"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestTooLongInvalidShortUrl(t *testing.T) {
	db := createDB()
	url := "dsfklaslkfjasldjfasdkfjasdlfjklasdjkfasdklflkasdjfldsfjasdshdjfaa"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestCharacterInvalidShortUrl(t *testing.T) {
	db := createDB()
	url := "/asdfsa=sdfasdf?q=klasd"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}
