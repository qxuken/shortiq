package db

type DB interface {
	GetLink(short string) (string, error)
	SetLink(url, short string) (int64, error)
}
