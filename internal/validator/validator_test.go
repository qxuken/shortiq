package validator_test

import (
	"net/url"
	"testing"

	"github.com/qxuken/short/internal"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/validator"
)

const MOCK_URL_STR = "https://short.iq"

var (
	MOCK_URL, _ = url.Parse(MOCK_URL_STR)
)

func createConf() *internal.Config {
	config := new(internal.Config)
	config.PublicUrl = *MOCK_URL
	return config
}

func createDB(t *testing.T) *db.SqliteDB {
	db, err := db.ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
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
	db := createDB(t)
	url := "qxuken"
	res := validator.ValidateShortHandle(db, url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestComplexValidShortUrl(t *testing.T) {
	db := createDB(t)
	url := "AaZz09qxuke23n-sadfashfsD_1sdasfdf"
	res := validator.ValidateShortHandle(db, url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestShortInvalidShortUrl(t *testing.T) {
	db := createDB(t)
	url := "qxuk"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestTooLongInvalidShortUrl(t *testing.T) {
	db := createDB(t)
	url := "dsfklaslkfjasldjfasdkfjasdlfjklasdjkfasdklflkasdjfldsfjasdshdjfaa"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestCharacterInvalidShortUrl(t *testing.T) {
	db := createDB(t)
	url := "/asdfsa=sdfasdf?q=klasd"
	res := validator.ValidateShortHandle(db, url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}
