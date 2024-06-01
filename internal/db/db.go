package db

type DB interface {
	GetLink(shortUrl string) (string, error)
	SetLink(redirectUrl, shortUrl string) error
	GetLinkAnalytics(shortUrl string) ([]AnalyticsItem, error)
	LogVisit(shortUrl string, v AnalyticsItem) error
}

type AnalyticsItem struct {
	Country string `db:"country"`
	Referer string `db:"referer"`
	Ip      string `db:"ip"`
	Ts      int64  `db:"ts"`
}
