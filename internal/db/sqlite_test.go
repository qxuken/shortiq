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
	kvs := [][2]string{{"testkey1", "testval1"}, {"UpercaseKey", "UpercaseVal"}}

	for _, kv := range kvs {
		_, err := db.SetLink(kv[0], kv[1])
		if err != nil {
			t.Fatal(err)
		}
		v, err := db.GetLink(kv[0])
		if err != nil {
			t.Fatal(err)
		}
		if v != kv[1] {
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
