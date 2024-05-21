package db

import (
	"github.com/hashicorp/golang-lru/v2"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS link (
		short STRING,
		url STRING
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_link_short ON link(short);
	`
	getLink = "SELECT url FROM link WHERE short = ? LIMIT 1;"
	setLink = "INSERT INTO link (url, short) VALUES (?, ?)"
)

type SqliteDB struct {
	db    *sqlx.DB
	cache *lru.Cache[string, string]
}

func (db *SqliteDB) GetLink(short string) (string, error) {
	cached, ok := db.cache.Get(short)
	if ok {
		return cached, nil
	}

	var url string
	err := db.db.Get(&url, getLink, short)
	if err != nil {
		return "", err
	}

	db.cache.Add(short, url)

	return url, nil
}

func (db *SqliteDB) SetLink(url, short string) error {
	_, err := db.db.Exec(setLink, url, short)
	return err
}

func ConnectSqlite3(path string) (*SqliteDB, error) {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)

	cache, err := lru.New[string, string](512)
	if err != nil {
		return nil, err
	}

	return &SqliteDB{db, cache}, nil
}
