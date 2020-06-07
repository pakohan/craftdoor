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
	MemberID  *int64 `json:"member_id" db:"member_id"`
	Secret    string `json:"secret" db:"secret"`
	AccessKey string `json:"access_key" db:"access_key"`
}

// List returns all entries from the table
func (m *Model) List(ctx context.Context) ([]Key, error) {
	res := []Key{}
	err := m.db.SelectContext(ctx, &res, queryList)
	if err != nil {
		return nil, err
	}
	return res, nil
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

// IsAccessAllowed returns whether the key has access to that door at the current time
func (m *Model) IsAccessAllowed(ctx context.Context, keyID string, doorID int64) (bool, error) {
	var res bool
	err := m.db.GetContext(ctx, &res, accessAllowed, keyID, doorID)
	return res, err
}

func (m *Model) AssignMember(ctx context.Context, keyID, memberID int64) error {
	mID := &memberID
	if memberID == 0 {
		mID = nil
	}
	_, err := m.db.ExecContext(ctx, queryAssign, mID, keyID)
	return err
}

const (
	queryCreate = `
INSERT INTO "main"."key"
( secret,  access_key)
VALUES
(:secret, :access_key)
ON CONFLICT DO NOTHING`
	queryList = `
SELECT "id"
	, "member_id"
	, "secret"
	, "access_key"
FROM "key"
ORDER BY "id"`
	queryAssign = `
UPDATE "key"
SET "member_id" = ?
WHERE "id" = ?`
	accessAllowed = `
SELECT COUNT(*) > 0
FROM key
JOIN member_role
	ON (key.member_id = member_role.member_id)
JOIN door_role
	ON (member_role.role_id = door_role.role_id)
WHERE key.secret = ?
AND door_role.door_id = ?
AND (member_role.expires_at > CURRENT_TIMESTAMP OR member_role.expires_at IS NULL)
AND (door_role.daytime_begin_seconds < strftime('%s',CURRENT_TIMESTAMP) - strftime('%s', DATE(CURRENT_TIMESTAMP)) OR door_role.daytime_begin_seconds IS NULL)
AND (door_role.daytime_end_seconds > strftime('%s',CURRENT_TIMESTAMP) - strftime('%s', DATE(CURRENT_TIMESTAMP)) OR door_role.daytime_end_seconds IS NULL)`
)
