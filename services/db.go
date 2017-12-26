package services

import (
	"github.com/go-pg/pg"
)

var (
	db      *pg.DB
	initial bool
)

type Config struct {
	URL string
}

// Init the db connection (and something else?)
func Init(url string) error {
	if initial {
		return nil
	}
	option, err := pg.ParseURL(url)
	if err != nil {
		return err
	}
	db = pg.Connect(option)
	return nil
}
