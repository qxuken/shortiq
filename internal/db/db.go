package db

type DB interface {
	GetLink(short string) (string, error)
	SetLink(url, short string) error
	GetVisits(short string) ([]LinkVisit, error)
	LogVisit(short string, v LinkVisit) error
}

type LinkVisit struct {
	Country string `db:"country"`
	Origin  string `db:"origin"`
	Ts      int64  `db:"ts"`
}
