package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
	"periph.io/x/periph/experimental/devices/mfrc522"
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

func (s *Service) InitKey(ctx context.Context) error {
	return s.r.InitKey([16]byte{1, 2, 3}, [16]byte{4, 5, 6}, mfrc522.Key{1, 2, 3, 4, 5, 6}, mfrc522.Key{1, 2, 3, 4, 5, 6}, mfrc522.Key{6, 5, 4, 3, 2, 1})
}
