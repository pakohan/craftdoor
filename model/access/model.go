package access

import "github.com/jmoiron/sqlx"

// Model accesses the db
type Model struct {
	db *sqlx.DB
}

// New returns a new model
func New(db *sqlx.DB) *Model {
	return &Model{db: db}
}
