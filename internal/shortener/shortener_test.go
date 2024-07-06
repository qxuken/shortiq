package shortener_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/shortener"
)

func createDB(t *testing.T) *db.SqliteDB {
	db, err := db.ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGeneratingUrls(t *testing.T) {
	short1 := shortener.ShortUrl()
	t.Logf("short1: %v\n", short1)
	if short1 == "" {
		t.Fatalf("Function did not generated hash")
	}
}

func TestGeneratingCollisionFree(t *testing.T) {
	db := createDB(t)
	urls := make([]string, 10_000)
	for i := 0; i < 5000; i++ {
		testS := fmt.Sprintf("test %d", i)
		urls[i] = testS
		urls[i+5000] = testS
	}

	for _, url := range urls {
		t.Logf("url: %v\n", url)
		short, err := shortener.ShortUrlChecked(db)
		t.Logf("short: %v\n", short)
		db.SetLink(url, short)
		t.Logf("err: %v\n", err)
		if err != nil {
			log.Fatal(err)
		}
	}
}
