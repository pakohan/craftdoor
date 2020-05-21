package service

import (
	"sync"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/config"
	"github.com/pakohan/craftdoor/lib"
	"github.com/pakohan/craftdoor/model"
)

type Service struct {
	m model.Model

	current      uuid.UUID
	lock         *sync.Mutex
	listeners    map[uuid.UUID]chan<- State
	currentState State
}

type State struct {
	s string
}

func New(cfg config.Config, m model.Model) (*Service, error) {
	s := &Service{
		m:         m,
		lock:      &sync.Mutex{},
		listeners: map[uuid.UUID]chan<- State{},
	}

	_, err := lib.NewReader(cfg, s)
	return s, err
}

func (s *Service) WaitForChange(id uuid.UUID) State {
	s.lock.Lock()
	if s.current == id {
		res := s.currentState
		s.lock.Unlock()
		return res
	}

	reqID := uuid.New()
	c := make(chan State)
	s.listeners[reqID] = c
	s.lock.Unlock()

	res := <-c

	s.lock.Lock()
	delete(s.listeners, reqID)
	s.lock.Unlock()

	return res
}

func (s *Service) Notify(state string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.current = uuid.New()
	s.currentState.s = state
	for _, l := range s.listeners {
		l <- s.currentState
		close(l)
	}
}
