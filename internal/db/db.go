package db

type MainDb interface {
	GetLink(shortUrl string) (string, error)
	SetLink(redirectUrl, shortUrl string) error
	GetLinks() ([]LinkItem, error)
	GetConfigItem(key string) (string, error)
	SetConfigItem(key, value string) error
}

type AuxiliaryDB interface {
	GetLinkAnalytics() ([]AnalyticsItem, error)
	LogVisit(v AnalyticsItem) error
	GetLinkStats(shortUrl string) (*LinkStats, error)
	GetAllLinksTrafficStats() ([]LinkTrafficStats, error)
	GetTrackingTotals() (*TrackingTotals, error)
	GetCountryStats(shortUrl string) ([]CountryStats, error)
	GetRefererStats(shortUrl string) ([]RefererStats, error)
	GetDailyClicks(shortUrl string, days int) ([]DailyStats, error)
	GetAllCountryStats() ([]CountryStats, error)
	GetAllRefererStats() ([]RefererStats, error)
	GetAllDailyClicks(days int) ([]DailyStats, error)
}

type LinkTrafficStats struct {
	ShortUrl       string `db:"short_url"`
	TotalClicks    int    `db:"total_clicks"`
	UniqueVisitors int    `db:"unique_visitors"`
}

type TrackingTotals struct {
	TotalClicks    int `db:"total_clicks"`
	UniqueVisitors int `db:"unique_visitors"`
}

type LinkItem struct {
	RedirectUrl string `db:"redirect_url" csv:"redirect_url"`
	ShortUrl    string `db:"short_url" csv:"short_url"`
}

type AnalyticsItem struct {
	ShortUrl string `db:"short_url" csv:"short_url"`
	Country  string `db:"country" csv:"country"`
	Referer  string `db:"referer" csv:"referer"`
	Ip       string `db:"ip" csv:"ip"`
	Ts       int64  `db:"ts" csv:"ts"`
}

type CountryStats struct {
	Country string `db:"country"`
	Count   int    `db:"count"`
}

type RefererStats struct {
	Referer string `db:"referer"`
	Count   int    `db:"count"`
}

type DailyStats struct {
	Date  string `db:"date"`
	Count int    `db:"count"`
}

type LinkStats struct {
	ShortUrl       string         `db:"short_url"`
	RedirectUrl    string         `db:"redirect_url"`
	TotalClicks    int            `db:"total_clicks"`
	UniqueVisitors int            `db:"unique_visitors"`
	TopCountries   []CountryStats `db:"-"`
	TopReferers    []RefererStats `db:"-"`
	DailyClicks    []DailyStats   `db:"-"`
}

type AllLinksStats struct {
	TotalLinks     int            `db:"total_links"`
	TotalClicks    int            `db:"total_clicks"`
	UniqueVisitors int            `db:"unique_visitors"`
	TopCountries   []CountryStats `db:"-"`
	TopReferers    []RefererStats `db:"-"`
	DailyClicks    []DailyStats   `db:"-"`
	LinksStats     []LinkStats    `db:"-"`
}
