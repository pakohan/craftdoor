package doorrole

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

// DoorRole represents a single row
type DoorRole struct {
	DoorID              int64  `json:"door_id" db:"door_id"`
	RoleID              int64  `json:"role_id" db:"role_id"`
	DaytimeBeginSeconds *int64 `json:"daytime_begin_seconds" db:"daytime_begin_seconds"`
	DaytimeEndSeconds   *int64 `json:"daytime_end_seconds" db:"daytime_end_seconds"`
}

// Create creates a new entry in the table
func (m *Model) Create(ctx context.Context, t *DoorRole) error {
	_, err := m.db.NamedExecContext(ctx, queryCreate, t)
	return err
}

// List returns all entries from the table
func (m *Model) List(ctx context.Context, doorID, roleID int64) ([]DoorRole, error) {
	res := []DoorRole{}
	err := m.db.SelectContext(ctx, &res, queryList, doorID, doorID, roleID, roleID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Update updates a single entry in the table
func (m *Model) Update(ctx context.Context, t DoorRole) error {
	_, err := m.db.NamedExecContext(ctx, queryUpdate, t)
	return err
}

// Delete deletes a single entry from the table
func (m *Model) Delete(ctx context.Context, doorID, roleID int64) error {
	_, err := m.db.ExecContext(ctx, queryDelete, doorID, roleID)
	return err
}

const (
	queryCreate = `
INSERT INTO "door_role"
("door_id", "role_id", "daytime_begin_seconds", "daytime_end_seconds")
VALUES
(:door_id,  :role_id,  :daytime_begin_seconds,  :daytime_end_seconds)
ON CONFLICT ("door_id", "role_id") DO UPDATE
SET "daytime_begin_seconds" = :daytime_begin_seconds
  , "daytime_end_seconds"   = :daytime_end_seconds`
	queryList = `
SELECT "door_id"
	, "role_id"
	, "daytime_begin_seconds"
	, "daytime_end_seconds"
FROM "door_role"
WHERE (? = 0 OR "door_id" = ?)
AND (? = 0 OR "role_id" = ?)`
	queryUpdate = `
UPDATE "door_role"
SET "daytime_begin_seconds" = :daytime_begin_seconds
  , "daytime_end_seconds"   = :daytime_end_seconds
WHERE "door_id" = :door_id
AND "role_id = :role_id`
	queryDelete = `
DELETE FROM "door_role"
WHERE door_id = ?
AND role_id = ?`
)
