package shortener_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/qxuken/short/internal/config"
	"github.com/qxuken/short/internal/db"
	"github.com/qxuken/short/internal/shortener"
)

func createDB() *db.SqliteDB {
	db := db.ConnectSqlite3(&config.Config{}, ":memory:")
	return db
}

func TestGeneratingUrls(t *testing.T) {
	short1 := shortener.ShortUrlWithLen(5)
	t.Logf("short1: %v\n", short1)
	if short1 == "" {
		t.Fatalf("Function did not generated hash")
	}
}

func TestGeneratingCollisionFree(t *testing.T) {
	db := createDB()
	urls := make([]string, 10_000)
	for i := 0; i < 5000; i++ {
		testS := fmt.Sprintf("test %d", i)
		urls[i] = testS
		urls[i+5000] = testS
	}

	for _, url := range urls {
		short, err := shortener.ShortUrlChecked(db, 5)
		db.SetLink(url, short)
		if err != nil {
			t.Logf("url: %v\n", url)
			t.Logf("short: %v\n", short)
			t.Logf("err: %v\n", err)
			log.Fatal(err)
		}
	}
}

var table = []struct {
	handleSize int
}{
	{5},
	{6},
	{7},
	{8},
}

func BenchmarkUrlShortener(b *testing.B) {
	for _, v := range table {
		b.Run(fmt.Sprintf("handle_size_%d", v.handleSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shortener.ShortUrlWithLen(v.handleSize)
			}
		})
	}
}
