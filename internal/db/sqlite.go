package db

import (
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
	db *sqlx.DB
}

func (db *SqliteDB) GetLink(short string) (url string, err error) {
	err = db.db.Get(&url, getLink, short)
	return
}

func (db *SqliteDB) SetLink(url, short string) (int64, error) {
	res, err := db.db.Exec(setLink, url, short)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ConnectSqlite3(path string) (*SqliteDB, error) {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)

	return &SqliteDB{db}, nil
}
