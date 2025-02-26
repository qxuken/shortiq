package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qxuken/short/internal/config"
)

const (
	schema = `
	PRAGMA journal_mode=WAL;
	PRAGMA foreign_keys = ON;
	PRAGMA synchronous = NORMAL;
	PRAGMA cache_size = 10000;
	PRAGMA temp_store = MEMORY;
	PRAGMA encoding = 'UTF-8';
	PRAGMA auto_vacuum = FULL;
	PRAGMA busy_timeout = 3000;
	PRAGMA optimize;
	PRAGMA mmap_size = 268435456;

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
	getLinks  = "SELECT redirect_url, short_url FROM link;"
	getVisits = "SELECT short_url, country, referer, ip, ts FROM analytics;"
	logVisit  = "INSERT INTO analytics (short_url, country, referer, ip, ts) VALUES (?, ?, ?, ?, ?);"
)

type SqliteDB struct {
	db *sqlx.DB
}

func ConnectSqlite3(conf *config.Config, path string) *SqliteDB {
	if conf.Verbose {
		log.Printf("Opening db on %v\n", path)
	}
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schema)
	return &SqliteDB{db}
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

func (db *SqliteDB) GetLinks() ([]LinkItem, error) {
	links := []LinkItem{}
	err := db.db.Select(&links, getLinks)
	return links, err
}

func (db *SqliteDB) GetLinkAnalytics() ([]AnalyticsItem, error) {
	visits := []AnalyticsItem{}
	err := db.db.Select(&visits, getVisits)
	return visits, err
}

func (db *SqliteDB) LogVisit(v AnalyticsItem) error {
	_, err := db.db.Exec(logVisit, v.ShortUrl, v.Country, v.Referer, v.Ip, v.Ts)
	return err
}
