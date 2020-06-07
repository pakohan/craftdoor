package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/pakohan/craftdoor/model/access"
	"github.com/pakohan/craftdoor/model/door"
	"github.com/pakohan/craftdoor/model/doorrole"
	"github.com/pakohan/craftdoor/model/key"
	"github.com/pakohan/craftdoor/model/member"
	"github.com/pakohan/craftdoor/model/memberrole"
	"github.com/pakohan/craftdoor/model/role"
)

// Model holds all models
type Model struct {
	AccessModel     *access.Model
	DoorModel       *door.Model
	DoorroleModel   *doorrole.Model
	KeyModel        *key.Model
	MemberModel     *member.Model
	MemberroleModel *memberrole.Model
	RoleModel       *role.Model
}

// New returns all models initialized
func New(db *sqlx.DB) Model {
	return Model{
		AccessModel:     access.New(db),
		DoorModel:       door.New(db),
		DoorroleModel:   doorrole.New(db),
		KeyModel:        key.New(db),
		MemberModel:     member.New(db),
		MemberroleModel: memberrole.New(db),
		RoleModel:       role.New(db),
	}
}
