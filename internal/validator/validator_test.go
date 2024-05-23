package validator_test

import (
	"testing"

	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/validator"
)

func createDB(t *testing.T) *db.SqliteDB {
	db, err := db.ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSimpleValidRedirectUrl(t *testing.T) {
	url := "https://github.com/qxuken"
	res := validator.ValidateRedirectUrl(url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestComplexValidRedirectUrl(t *testing.T) {
	url := "https://github.com/qxuken/key-val?query=some+hey"
	res := validator.ValidateRedirectUrl(url)
	if res != nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestProtocollesInvalidRedirectUrl(t *testing.T) {
	url := "github.com/qxuken/key-val?query=some+hey"
	res := validator.ValidateRedirectUrl(url)
	if res == nil {
		t.Fatalf("Failed to validate %s", url)
	}
}

func TestShortInvalidRedirectUrl(t *testing.T) {
	url := "/qxuken"
	res := validator.ValidateRedirectUrl(url)
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
