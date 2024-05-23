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
	
	CREATE TABLE IF NOT EXISTS link_visit (
		short STRING,
		country STRING,
		origin STRING,
		ts INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_link_visit_short ON link_visit(short);
	`
	getLink   = "SELECT url FROM link WHERE short = ? LIMIT 1;"
	setLink   = "INSERT INTO link (url, short) VALUES (?, ?)"
	getVisits = "SELECT country, origin, ts FROM link_visit WHERE short = ?"
	logVisit  = "INSERT INTO link_visit (short, country, origin, ts) VALUES (?, ?, ?, ?)"
)

type SqliteDB struct {
	db *sqlx.DB
}

func ConnectSqlite3(path string) (*SqliteDB, error) {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)

	return &SqliteDB{db}, nil
}

func (db *SqliteDB) GetLink(short string) (string, error) {
	var url string
	err := db.db.Get(&url, getLink, short)
	return url, err
}

func (db *SqliteDB) SetLink(url, short string) error {
	_, err := db.db.Exec(setLink, url, short)
	return err
}

func (db *SqliteDB) GetVisits(short string) ([]LinkVisit, error) {
	visits := []LinkVisit{}
	err := db.db.Select(&visits, getVisits, short)
	return visits, err
}

func (db *SqliteDB) LogVisit(short string, v LinkVisit) error {
	_, err := db.db.Exec(logVisit, short, v.Country, v.Origin, v.Ts)
	return err
}
