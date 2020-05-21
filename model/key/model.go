package key

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Model accesses the key table
type Model struct {
	db *sqlx.DB
}

// New returns a new model
func New(db *sqlx.DB) *Model {
	return &Model{db: db}
}

// Key represents a single row
type Key struct {
	ID        int64  `json:"id" db:"id"`
	Secret    string `json:"secret" db:"secret"`
	AccessKey string `json:"access_key" db:"access_key"`
}

// Create inserts a new row into the table
func (m *Model) Create(ctx context.Context, k *Key) error {
	res, err := m.db.NamedExecContext(ctx, queryCreate, k)
	if err != nil {
		return err
	}
	k.ID, err = res.LastInsertId()
	return err
}

const (
	queryCreate = `
INSERT INTO "main"."key"
( secret,  access_key)
VALUES
(:secret, :access_key`
)
