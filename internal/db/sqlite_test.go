package db_test

import (
	"slices"
	"testing"

	"github.com/qxuken/short/internal/db"
)

func createDB(t *testing.T) *db.SqliteDB {
	db, err := db.ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSqlite3KVs(t *testing.T) {
	db := createDB(t)
	kvs := [][2]string{{"testurl", "testshort"}, {"UpercaseUrl", "UpercaseSHORT"}}

	for _, kv := range kvs {
		err := db.SetLink(kv[0], kv[1])
		if err != nil {
			t.Fatal(err)
		}
		v, err := db.GetLink(kv[1])
		if err != nil {
			t.Fatal(err)
		}
		if v != kv[0] {
			t.Fatalf("Value missmatch for %s, expected %s but found %s", kv[0], kv[1], v)
		}
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
	tv := []db.AnalyticsItem{{"c1", "o1", "192.168.0.1", 1}, {"c2", "o2", "192.168.0.2", 2}}
	db := createDB(t)
	for _, v := range tv {
		err := db.LogVisit("s", v)
		if err != nil {
			t.Fatal(err)
		}
	}
	r, err := db.GetLinkAnalytics("s")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(tv, r) {
		t.Fatalf("Value missmatch, got response (len = %v, val = %v)", len(r), r)
	}
}
