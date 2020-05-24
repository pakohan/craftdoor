package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
	"periph.io/x/periph/experimental/devices/mfrc522"
)

// Service contains the business logic
type Service struct {
	m  model.Model
	r  lib.Reader
	cl *lib.ChangeListener
}

// New returns a new service instance
func New(m model.Model, r lib.Reader, cl *lib.ChangeListener) *Service {
	return &Service{
		m:  m,
		r:  r,
		cl: cl,
	}
}

// WaitForChange returns as soon as the state id is different to the one passed
func (s *Service) WaitForChange(ctx context.Context, id uuid.UUID) (lib.State, error) {
	return s.cl.WaitForChange(ctx, id)
}

// InitKey writes a sector on a RFID card
func (s *Service) InitKey(ctx context.Context) error {
	return s.r.InitKey([16]byte{1, 2, 3}, [16]byte{4, 5, 6}, mfrc522.DefaultKey, mfrc522.Key{1, 2, 3, 4, 5, 6}, mfrc522.Key{6, 5, 4, 3, 2, 1})
}
