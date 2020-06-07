package lib

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

// ChangeListener informs any registered listener about state changes of the reader.
// Every state change gets a new UUID assigned in order to let the listeners track whether there were any changes since the last time the state was returned
type ChangeListener struct {
	Reader
	lock         *sync.Mutex
	listeners    map[uuid.UUID]chan<- State
	currentState State
}

// NewChangeListener returns a new ChangeListener ready to be used
func NewChangeListener() *ChangeListener {
	return &ChangeListener{
		lock:      &sync.Mutex{},
		listeners: map[uuid.UUID]chan<- State{},
	}
}

// WaitForChange checks whether the id equals the one of the current state. If not, the state is immediately returned.
// If the id is the same it means the listener knew about the current state already and the function returns either on timeout
// or when the state changes
func (cl *ChangeListener) WaitForChange(ctx context.Context, id uuid.UUID) (State, error) {
	cl.lock.Lock()
	if cl.currentState.ID != id {
		res := cl.currentState
		cl.lock.Unlock()
		return res, nil
	}

	reqID := uuid.New()
	c := make(chan State)
	cl.listeners[reqID] = c
	cl.lock.Unlock()

	defer func() {
		cl.lock.Lock()
		delete(cl.listeners, reqID)
		cl.lock.Unlock()
	}()

	select {
	case res := <-c:
		return res, nil
	case <-ctx.Done():
		return State{}, ctx.Err()
	}
}

func (cl *ChangeListener) ReturnFirstKey(ctx context.Context) (State, error) {
	cl.lock.Lock()
	if cl.currentState.IsCardAvailable {
		res := cl.currentState
		cl.lock.Unlock()
		return res, nil
	}

	reqID := uuid.New()
	c := make(chan State)
	cl.listeners[reqID] = c
	cl.lock.Unlock()

	defer func() {
		cl.lock.Lock()
		delete(cl.listeners, reqID)
		cl.lock.Unlock()
	}()

	select {
	case res := <-c:
		return res, nil
	case <-ctx.Done():
		return State{}, ctx.Err()
	}
}

// Notify changes the current state
func (cl *ChangeListener) Notify(data [3]string) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	log.Printf("state changed to: %s", data[0])

	cl.currentState.IsCardAvailable = data[0] != ""
	cl.currentState.ID = uuid.New()
	cl.currentState.CardData = make([]string, len(data))
	for i, s := range data {
		cl.currentState.CardData[i] = fmt.Sprintf("% x", s)
	}

	for _, l := range cl.listeners {
		l <- cl.currentState
		close(l)
	}
}
