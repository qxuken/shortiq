package db

type DB interface {
	GetLink(alias string) (string, error)
	SetLink(alias, url string) (int64, error)
}
