package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
)

type Service struct {
	m  model.Model
	r  *lib.Reader
	cl *lib.ChangeListener
}

func New(m model.Model, r *lib.Reader, cl *lib.ChangeListener) *Service {
	return &Service{
		m:  m,
		r:  r,
		cl: cl,
	}
}

func (s *Service) WaitForChange(ctx context.Context, id uuid.UUID) (lib.State, error) {
	return s.cl.WaitForChange(ctx, id)
}
