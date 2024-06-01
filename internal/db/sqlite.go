package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS link (
		redirect_url STRING,
		short_url STRING
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_link_short ON link(short_url);
	
	CREATE TABLE IF NOT EXISTS analytics (
		short_url STRING,
		country STRING,
		referer STRING,
		ip STRING,
		ts INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_analytics_short ON analytics(short_url);
	`
	getLink   = "SELECT redirect_url FROM link WHERE short_url = ? LIMIT 1;"
	setLink   = "INSERT INTO link (redirect_url, short_url) VALUES (?, ?);"
	getVisits = "SELECT country, referer, ip, ts FROM analytics WHERE short_url = ?;"
	logVisit  = "INSERT INTO analytics (short_url, country, referer, ip, ts) VALUES (?, ?, ?, ?, ?);"
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

func (db *SqliteDB) GetLink(shortUrl string) (string, error) {
	var url string
	err := db.db.Get(&url, getLink, shortUrl)
	return url, err
}

func (db *SqliteDB) SetLink(redirectUrl, shortUrl string) error {
	_, err := db.db.Exec(setLink, redirectUrl, shortUrl)
	return err
}

func (db *SqliteDB) GetLinkAnalytics(shortUrl string) ([]AnalyticsItem, error) {
	visits := []AnalyticsItem{}
	err := db.db.Select(&visits, getVisits, shortUrl)
	return visits, err
}

func (db *SqliteDB) LogVisit(shortUrl string, v AnalyticsItem) error {
	_, err := db.db.Exec(logVisit, shortUrl, v.Country, v.Referer, v.Ip, v.Ts)
	return err
}
