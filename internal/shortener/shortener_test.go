package shortener_test

import (
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
	short1, short2 := shortener.ShortUrl("test1"), shortener.ShortUrl("test2")
	t.Logf("short1: %v\n", short1)
	t.Logf("short2: %v\n", short2)
	if short1 == short2 {
		t.Fatalf("Function generated equivalent hash, %s", short1)
	}
}

func TestGeneratingColissionFree(t *testing.T) {
	db := createDB(t)
	urls := []string{"test1", "test2", "test3", "test1", "aaa", "aaa"}

	for _, url := range urls {
		t.Logf("url: %v\n", url)
		short, err := shortener.ShortUrlChecked(db, url)
		t.Logf("short: %v\n", short)
		db.SetLink(url, short)
		t.Logf("err: %v\n", err)
		if err != nil {
			log.Fatal(err)
		}
	}
}
