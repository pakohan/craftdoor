package role

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Model accesses the db
type Model struct {
	db *sqlx.DB
}

// New returns a new model
func New(db *sqlx.DB) *Model {
	return &Model{db: db}
}

// Role represents a single row
type Role struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Create creates a new entry in the table
func (m *Model) Create(ctx context.Context, t *Role) error {
	res, err := m.db.NamedExecContext(ctx, queryCreate, t)
	if err != nil {
		return err
	}
	t.ID, err = res.LastInsertId()
	return err
}

// List returns all entries from the table
func (m *Model) List(ctx context.Context) ([]Role, error) {
	res := []Role{}
	err := m.db.SelectContext(ctx, &res, queryList)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Update updates a single entry in the table
func (m *Model) Update(ctx context.Context, t Role) error {
	_, err := m.db.NamedExecContext(ctx, queryUpdate, t)
	return err
}

// Delete deletes a single entry from the table
func (m *Model) Delete(ctx context.Context, id int64) error {
	_, err := m.db.ExecContext(ctx, queryDelete, id)
	return err
}

const (
	queryCreate = `
INSERT INTO "role"
("name")
VALUES
(:name)`
	queryList = `
SELECT "id"
	, "name"
FROM "role"
ORDER BY "id"`
	queryUpdate = `
UPDATE "role"
SET "name" = :name
WHERE "id" = :id`
	queryDelete = `
DELETE FROM "role"
WHERE id = ?`
)
