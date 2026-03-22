package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qxuken/short/internal/config"
)

type UrlStore struct {
	db *sqlx.DB
}

func ConnectUrlStore(conf *config.Config, path string) *UrlStore {
	if conf.Verbose {
		log.Printf("Opening url store on %v\n", path)
	}
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(`
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
	
	CREATE TABLE IF NOT EXISTS app_config (
		key STRING,
		value STRING
	);
	CREATE INDEX IF NOT EXISTS idx_app_config_key ON app_config(key);
	`)
	return &UrlStore{db}
}

func (s *UrlStore) GetLink(shortUrl string) (url string, err error) {
	err = s.db.Get(&url, "SELECT redirect_url FROM link WHERE short_url = ? LIMIT 1;", shortUrl)
	return
}

func (s *UrlStore) SetLink(redirectUrl, shortUrl string) (err error) {
	_, err = s.db.Exec("INSERT INTO link (redirect_url, short_url) VALUES (?, ?);", redirectUrl, shortUrl)
	return
}

func (s *UrlStore) GetLinks() (links []LinkItem, err error) {
	err = s.db.Select(&links, "SELECT redirect_url, short_url FROM link;")
	return
}

func (s *UrlStore) GetConfigItem(key string) (value string, err error) {
	err = s.db.Get(&value, "SELECT value FROM app_config WHERE key = ?;", key)
	return
}

func (s *UrlStore) SetConfigItem(key, value string) (err error) {
	_, err = s.db.Exec("INSERT INTO app_config (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = excluded.value", key, value)
	return
}

type TrackingStore struct {
	db *sqlx.DB
}

func ConnectTrackingStore(conf *config.Config, path string) *TrackingStore {
	if conf.Verbose {
		log.Printf("Opening tracking store on %v\n", path)
	}
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(`
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

	CREATE TABLE IF NOT EXISTS analytics (
		short_url STRING,
		country STRING,
		referer STRING,
		ip STRING,
		ts INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_analytics_short ON analytics(short_url);
	`)
	return &TrackingStore{db}
}

func (s *TrackingStore) GetLinkAnalytics() (visits []AnalyticsItem, err error) {
	err = s.db.Select(&visits, "SELECT short_url, country, referer, ip, ts FROM analytics;")
	return
}

func (s *TrackingStore) LogVisit(v AnalyticsItem) (err error) {
	_, err = s.db.Exec("INSERT INTO analytics (short_url, country, referer, ip, ts) VALUES (?, ?, ?, ?, ?);", v.ShortUrl, v.Country, v.Referer, v.Ip, v.Ts)
	return
}

func (s *TrackingStore) GetLinkStats(shortUrl string) (*LinkStats, error) {
	query := `
		SELECT 
			COALESCE(COUNT(*), 0) as total_clicks,
			COALESCE(COUNT(DISTINCT ip), 0) as unique_visitors
		FROM analytics WHERE short_url = ?
	`
	var stats struct {
		TotalClicks    int `db:"total_clicks"`
		UniqueVisitors int `db:"unique_visitors"`
	}
	err := s.db.Get(&stats, query, shortUrl)
	if err != nil {
		return nil, err
	}
	return &LinkStats{
		ShortUrl:       shortUrl,
		TotalClicks:    stats.TotalClicks,
		UniqueVisitors: stats.UniqueVisitors,
	}, nil
}

func (s *TrackingStore) GetAllLinksTrafficStats() ([]LinkTrafficStats, error) {
	query := `
		SELECT
			short_url,
			COALESCE(COUNT(*), 0) as total_clicks,
			COALESCE(COUNT(DISTINCT ip), 0) as unique_visitors
		FROM analytics
		GROUP BY short_url
		ORDER BY total_clicks DESC
	`
	var stats []LinkTrafficStats
	err := s.db.Select(&stats, query)
	return stats, err
}

func (s *TrackingStore) GetTrackingTotals() (*TrackingTotals, error) {
	query := `
		SELECT
			COALESCE(COUNT(*), 0) as total_clicks,
			COALESCE(COUNT(DISTINCT ip), 0) as unique_visitors
		FROM analytics
	`
	var stats TrackingTotals
	err := s.db.Get(&stats, query)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (s *TrackingStore) GetCountryStats(shortUrl string) ([]CountryStats, error) {
	query := `
		SELECT country, COUNT(*) as count 
		FROM analytics 
		WHERE short_url = ? AND country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT 5
	`
	var stats []CountryStats
	err := s.db.Select(&stats, query, shortUrl)
	return stats, err
}

func (s *TrackingStore) GetRefererStats(shortUrl string) ([]RefererStats, error) {
	query := `
		SELECT referer, COUNT(*) as count 
		FROM analytics 
		WHERE short_url = ? AND referer != ''
		GROUP BY referer 
		ORDER BY count DESC 
		LIMIT 5
	`
	var stats []RefererStats
	err := s.db.Select(&stats, query, shortUrl)
	return stats, err
}

func (s *TrackingStore) GetDailyClicks(shortUrl string, days int) ([]DailyStats, error) {
	query := `
		SELECT 
			date(ts, 'unixepoch') as date,
			COUNT(*) as count
		FROM analytics
		WHERE short_url = ?
		GROUP BY date
		ORDER BY date DESC
		LIMIT ?
	`
	var stats []DailyStats
	err := s.db.Select(&stats, query, shortUrl, days)
	return stats, err
}

func (s *TrackingStore) GetAllCountryStats() ([]CountryStats, error) {
	query := `
		SELECT country, COUNT(*) as count 
		FROM analytics 
		WHERE country != ''
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT 10
	`
	var stats []CountryStats
	err := s.db.Select(&stats, query)
	return stats, err
}

func (s *TrackingStore) GetAllRefererStats() ([]RefererStats, error) {
	query := `
		SELECT referer, COUNT(*) as count 
		FROM analytics 
		WHERE referer != ''
		GROUP BY referer 
		ORDER BY count DESC 
		LIMIT 10
	`
	var stats []RefererStats
	err := s.db.Select(&stats, query)
	return stats, err
}

func (s *TrackingStore) GetAllDailyClicks(days int) ([]DailyStats, error) {
	query := `
		SELECT 
			date(ts, 'unixepoch') as date,
			COUNT(*) as count
		FROM analytics
		GROUP BY date
		ORDER BY date DESC
		LIMIT ?
	`
	var stats []DailyStats
	err := s.db.Select(&stats, query, days)
	return stats, err
}
