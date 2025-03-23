package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(user string, pass string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432?sslmode=disable",
		user, pass)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
