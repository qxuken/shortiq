package db

type DB interface {
	GetLink(shortUrl string) (string, error)
	SetLink(redirectUrl, shortUrl string) error
	GetLinks() ([]LinkItem, error)
	GetLinkAnalytics() ([]AnalyticsItem, error)
	LogVisit(v AnalyticsItem) error
	GetConfigKey(key string) (string, error)
	SetConfigKey(key, value string) error
}

type LinkItem struct {
	RedirectUrl string `db:"redirect_url" csv:"redirect_url"`
	ShortUrl    string `db:"short_url" csv:"short_url"`
}

type AnalyticsItem struct {
	ShortUrl string `db:"short_url" csv:"short_url"`
	Country  string `db:"country" csv:"country"`
	Referer  string `db:"referer" csv:"refere"`
	Ip       string `db:"ip" csv:"ip"`
	Ts       int64  `db:"ts" csv:"ts"`
}
