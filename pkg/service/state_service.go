package service

import (
	"sync"
	"time"

	"github.com/ourorg/goui/pkg/domain"
)

type DefaultStateStore struct {
	mu     sync.RWMutex
	curr   *domain.State
	subs   []func(*domain.State)
	stateReg *StateRegistry
}

func NewDefaultStateStore(reg *StateRegistry) *DefaultStateStore {
	return &DefaultStateStore{stateReg: reg}
}

func (s *DefaultStateStore) Init(initialID int) error {
	if s.stateReg == nil { return nil }
	stIdx := s.stateReg.Index()
	if len(stIdx) == 0 { return nil }
	curr := stIdx[initialID]
	s.Commit(func(_ *domain.State) (*domain.State, bool) { return &curr, true })
	return nil
}

func (s *DefaultStateStore) Current() *domain.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.curr
}

func (s *DefaultStateStore) Commit(f func(*domain.State) (*domain.State, bool)) (*domain.State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	next, ok := f(s.curr)
	if ok {
		s.curr = next
		for _, fn := range s.subs {
			if fn != nil { fn(s.curr) }
		}
	}
	return s.curr, ok
}

func (s *DefaultStateStore) Subscribe(fn func(*domain.State)) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs = append(s.subs, fn)
	idx := len(s.subs) - 1
	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if idx < len(s.subs) {
			s.subs[idx] = nil
		}
	}
}

type StateService struct {
	store    StateStore
	history  StateHistory
	stateReg *StateRegistry
}

func NewStateService(store StateStore, reg *StateRegistry) *StateService {
	return &StateService{store: store, stateReg: reg}
}

func (s *StateService) Init(initialID int) error {
	// set first state as current
	if s.stateReg == nil { return nil }
	stIdx := s.stateReg.Index()
	if len(stIdx) == 0 { return nil }
	curr := stIdx[initialID]
	s.store.Commit(func(_ *domain.State) (*domain.State, bool) { return &curr, true })
	return nil
}

func (s *StateService) Current() *domain.State {
	return s.store.Current()
}

func (s *StateService) SetNextState(toID int, mutateArgs func(map[string]interface{})) error {
	curr := s.store.Current()
	fromID := -1
	if curr != nil { fromID = curr.ID }

	stateMap := s.stateReg.Index()
	next, ok := stateMap[toID]
	if !ok {
		// Stay in current state if target doesn't exist
		return nil
	}

	// copy to avoid mutating registry copy
	cp := next
	if cp.Args == nil { cp.Args = map[string]interface{}{} }
	if mutateArgs != nil { mutateArgs(cp.Args) }

	s.store.Commit(func(_ *domain.State) (*domain.State, bool) { return &cp, true })

	// Record transition for undo/redo
	s.history.Undo = append(s.history.Undo, Transition{
		FromID: fromID,
		ToID: toID,
		Cause: "SetNextState",
		At: time.Now(),
	})
	// Clear redo stack when new action occurs
	s.history.Redo = nil

	return nil
}

func (s *StateService) History() StateHistory {
	return s.history
}

func (s *StateService) Undo() bool {
	if len(s.history.Undo) == 0 { return false }

	// Pop from undo stack
	lastTransition := s.history.Undo[len(s.history.Undo)-1]
	s.history.Undo = s.history.Undo[:len(s.history.Undo)-1]

	// Move the undone transition to redo stack (preserving original direction for redo)
	s.history.Redo = append(s.history.Redo, lastTransition)

	// Go to previous state
	if lastTransition.FromID >= 0 {
		stateMap := s.stateReg.Index()
		if prevState, ok := stateMap[lastTransition.FromID]; ok {
			s.store.Commit(func(_ *domain.State) (*domain.State, bool) {
				return &prevState, true
			})
		}
	}

	return true
}

func (s *StateService) Redo() bool {
	if len(s.history.Redo) == 0 { return false }

	// Pop from redo stack
	lastTransition := s.history.Redo[len(s.history.Redo)-1]
	s.history.Redo = s.history.Redo[:len(s.history.Redo)-1]

	// Move the redone transition back to undo stack
	s.history.Undo = append(s.history.Undo, lastTransition)

	// Go to the target state (the state we're redoing to)
	stateMap := s.stateReg.Index()
	if nextState, ok := stateMap[lastTransition.ToID]; ok {
		s.store.Commit(func(_ *domain.State) (*domain.State, bool) {
			return &nextState, true
		})
	}

	return true
}

func (s *StateService) Push(id int) error {
	return s.SetNextState(id, nil)
}

func (s *StateService) Pop() error {
	return nil
}