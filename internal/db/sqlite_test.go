package db_test

import (
	"slices"
	"testing"

	mdb "github.com/qxuken/short/internal/db"
)

func createDB(t *testing.T) *mdb.SqliteDB {
	db, err := mdb.ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSqlite3KVs(t *testing.T) {
	db := createDB(t)
	kvs := []mdb.LinkItem{{"testurl", "testshort"}, {"UpercaseUrl", "UpercaseSHORT"}}

	for _, kv := range kvs {
		err := db.SetLink(kv.RedirectUrl, kv.ShortUrl)
		if err != nil {
			t.Fatal(err)
		}
		v, err := db.GetLink(kv.ShortUrl)
		if err != nil {
			t.Fatal(err)
		}
		if v != kv.RedirectUrl {
			t.Fatalf("Value missmatch for %s, expected %s but found %s", kv.ShortUrl, kv.RedirectUrl, v)
		}
	}

	r, err := db.GetLinks()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(kvs, r) {
		t.Fatalf("Value missmatch, got response (len = %v, val = %v)", len(r), r)
	}
}

func TestSqlite3EmptyKey(t *testing.T) {
	db := createDB(t)
	v, err := db.GetLink("empty")
	if err == nil || v != "" {
		t.Fatal("Found value where it shouldnt be")
	}
}

func TestSqlite3LogVisit(t *testing.T) {
	tv := []mdb.AnalyticsItem{{"s", "c1", "o1", "192.168.0.1", 1}, {"s", "c2", "o2", "192.168.0.2", 2}, {"s2", "c3", "o3", "192.168.0.2", 3}}
	db := createDB(t)
	for _, v := range tv {
		err := db.LogVisit(v)
		if err != nil {
			t.Fatal(err)
		}
	}
	r, err := db.GetLinkAnalytics()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tv, r) {
		t.Fatalf("Value missmatch, got response (len = %v, val = %v)", len(r), r)
	}
}
