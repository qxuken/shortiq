package db

import "testing"

func createDB(t *testing.T) *SqliteDB {
	db, err := ConnectSqlite3(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSqlite3KVs(t *testing.T) {
	db := createDB(t)
	kvs := [][2]string{{"testurl", "testshort"}, {"UpercaseUrl", "UpercaseSHORT"}}

	for _, kv := range kvs {
		_, err := db.SetLink(kv[0], kv[1])
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
