package service

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/pakohan/craftdoor/config"
	"github.com/pakohan/craftdoor/model"
)

type Service struct {
	m model.Model

	lock         *sync.Mutex
	listeners    map[uuid.UUID]chan<- State
	currentState State
}

type State struct {
	S  string
	ID uuid.UUID
}

func New(cfg config.Config, m model.Model) *Service {
	return &Service{
		m:         m,
		lock:      &sync.Mutex{},
		listeners: map[uuid.UUID]chan<- State{},
	}
}

func (s *Service) WaitForChange(id uuid.UUID) State {
	s.lock.Lock()
	if s.currentState.ID == id {
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

	log.Printf("state changed to '%s'", state)

	s.currentState.S = state
	s.currentState.ID = uuid.New()
	for _, l := range s.listeners {
		l <- s.currentState
		close(l)
	}
}
