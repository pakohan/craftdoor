package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/pakohan/craftdoor/model/access"
	"github.com/pakohan/craftdoor/model/key"
)

type Model struct {
	KeyModel    *key.Model
	AccessModel *access.Model
}

func New(db *sqlx.DB) Model {
	return Model{
		KeyModel:    key.New(db),
		AccessModel: access.New(db),
	}
}
