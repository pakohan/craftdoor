package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
	"github.com/pakohan/craftdoor/model/key"
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
	s := &Service{
		m:  m,
		r:  r,
		cl: cl,
	}
	go s.door()
	return s
}

// WaitForChange returns as soon as the state id is different to the one passed
func (s *Service) WaitForChange(ctx context.Context, id uuid.UUID) (lib.State, error) {
	return s.cl.WaitForChange(ctx, id)
}

// InitKey writes a sector on a RFID card
func (s *Service) InitKey(ctx context.Context) error {
	return s.r.InitKey([16]byte{1, 2, 3}, [16]byte{4, 5, 6}, mfrc522.DefaultKey, mfrc522.Key{1, 2, 3, 4, 5, 6}, mfrc522.Key{6, 5, 4, 3, 2, 1})
}

func (s *Service) RegisterKey(ctx context.Context) (key.Key, error) {
	log.Printf("registering key")
	state, err := s.cl.ReturnFirstKey(ctx)
	if err != nil {
		return key.Key{}, err
	}

	log.Printf("got key %s", state.CardData[0])

	k := key.Key{
		Secret:    state.CardData[0],
		AccessKey: uuid.New().String(),
	}
	err = s.m.KeyModel.Create(ctx, &k)
	if err != nil {
		return key.Key{}, err
	}

	log.Printf("inserted as key id %d", k.ID)

	return k, nil
}

func (s *Service) door() {
	uid := uuid.Nil
	for {
		state, err := s.WaitForChange(context.Background(), uid)
		if err != nil {
			log.Printf("got err waiting for change: %s", err.Error())
			continue
		} else if !state.IsCardAvailable {
			continue
		}
		uid = state.ID

		res, err := s.m.KeyModel.IsAccessAllowed(context.Background(), state.CardData[0], 1)
		if err != nil {
			log.Printf("got err checking whether key is allowed to access: %s", err.Error)
		} else {
			log.Printf("key %s is allowed to access door %d -> %t", state.CardData[0], 1, res)
		}
	}
}
