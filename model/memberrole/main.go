package memberrole

import (
	"context"
	"time"

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

// MemberRole represents a single row
type MemberRole struct {
	MemberID  int64      `json:"member_id" db:"member_id"`
	RoleID    int64      `json:"role_id" db:"role_id"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
}

// Create creates a new entry in the table
func (m *Model) Create(ctx context.Context, t *MemberRole) error {
	_, err := m.db.NamedExecContext(ctx, queryCreate, t)
	return err
}

// List returns all entries from the table
func (m *Model) List(ctx context.Context, memberID, roleID uint64) ([]MemberRole, error) {
	res := []MemberRole{}
	err := m.db.SelectContext(ctx, &res, queryList, memberID, memberID, roleID, roleID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Update updates a single entry in the table
func (m *Model) Update(ctx context.Context, t MemberRole) error {
	_, err := m.db.NamedExecContext(ctx, queryUpdate, t)
	return err
}

// Delete deletes a single entry from the table
func (m *Model) Delete(ctx context.Context, memberID, roleID int64) error {
	_, err := m.db.ExecContext(ctx, queryDelete, memberID, roleID)
	return err
}

const (
	queryCreate = `
INSERT INTO "member_role"
("member_id", "role_id", "expires_at")
VALUES
(:member_id,  :role_id,  :expires_at)
ON CONFLICT ("member_id", "role_id") DO UPDATE
SET "expires_at" = :expires_at`
	queryList = `
SELECT "member_id"
	, "role_id"
	, "expires_at"
FROM "member_role"
WHERE (? = 0 OR "member_id" = ?)
AND (? = 0 OR "role_id" = ?)`
	queryUpdate = `
UPDATE "member_role"
SET "expires_at" = :expires_at
WHERE "member_id" = :member_id
AND "role_id = :role_id`
	queryDelete = `
DELETE FROM "member_role"
WHERE member_id = ?
AND role_id = ?`
)
