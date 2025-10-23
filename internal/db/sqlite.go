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
		short_url STRING,
		created_at STRING
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

	CREATE TABLE IF NOT EXISTS app_config (
		key STRING,
		value STRING
	);
	CREATE INDEX IF NOT EXISTS idx_app_config_key ON app_config(key);
	`
	queryLink    = "SELECT redirect_url FROM link WHERE short_url = ? LIMIT 1;"
	insertLink   = "INSERT INTO link (redirect_url, short_url) VALUES (?, ?);"
	queryLinks   = "SELECT redirect_url, short_url FROM link;"
	queryVisits  = "SELECT short_url, country, referer, ip, ts FROM analytics;"
	insertVisit  = "INSERT INTO analytics (short_url, country, referer, ip, ts) VALUES (?, ?, ?, ?, ?);"
	queryConfig  = "SELECT value FROM app_config WHERE key = ?;"
	upsertConfig = "INSERT INTO app_config (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = excluded.value"
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

func (db *SqliteDB) GetLink(shortUrl string) (url string, err error) {
	err = db.db.Get(&url, queryLink, shortUrl)
	return
}

func (db *SqliteDB) SetLink(redirectUrl, shortUrl string) (err error) {
	_, err = db.db.Exec(insertLink, redirectUrl, shortUrl)
	return
}

func (db *SqliteDB) GetLinks() (links []LinkItem, err error) {
	err = db.db.Select(&links, queryLinks)
	return
}

func (db *SqliteDB) GetLinkAnalytics() (visits []AnalyticsItem, err error) {
	err = db.db.Select(&visits, queryVisits)
	return
}

func (db *SqliteDB) LogVisit(v AnalyticsItem) (err error) {
	_, err = db.db.Exec(insertVisit, v.ShortUrl, v.Country, v.Referer, v.Ip, v.Ts)
	return
}

func (db *SqliteDB) GetConfigItem(key string) (value string, err error) {
	err = db.db.Get(&value, queryConfig, key)
	return
}

func (db *SqliteDB) SetConfigItem(key, value string) (err error) {
	_, err = db.db.Exec(upsertConfig, key, value)
	return
}
