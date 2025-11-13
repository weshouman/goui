package service

import (
	"time"

	"github.com/ourorg/goui/pkg/domain"
	"github.com/ourorg/goui/pkg/spec"
)

// SpecProvider builds a renderer-agnostic spec from a state.
type SpecProvider interface {
	BuildSpec(*domain.State) spec.Spec
	ApplyFilter(spec.Spec, map[string]interface{}) spec.Spec
}

// StateProvider is the single source of truth for state and transitions.
type StateProvider interface {
	Init(initialID int) error
	Current() *domain.State
	SetNextState(toID int, mutateArgs func(map[string]interface{})) error
	History() StateHistory
	Undo() bool
	Redo() bool
	Push(id int) error
	Pop() error
}

// StateWriter is the write-only subset Engine exposes to UIs and commands.
type StateWriter interface {
	SetNextState(id int, mutateArgs func(map[string]interface{})) error
	Push(id int) error
	Pop() error
}

// StateStore is the low-level storage for StateService.
type StateStore interface {
	Init(initialID int) error
	Current() *domain.State
	Commit(f func(*domain.State) (*domain.State, bool)) (*domain.State, bool)
	Subscribe(fn func(*domain.State)) func()
}

// ModeProvider owns current mode.
type ModeProvider interface {
	SetMode(int)
	CurrentMode() int
}

// CommandProvider handles suggestions, autocomplete, and dispatch helpers.
type CommandProvider interface {
	Suggestions(prefix string) []string
	Autocomplete(prefix string) []string
	TouchHistory(cmd string)
	Resolve(alias string) (*domain.Command, bool)
	Dispatch(alias string, args []string) (string, error)
}

// Simple history structs to mirror the diagram.
type Transition struct {
	FromID int
	ToID   int
	Cause  string
	At     time.Time
}

type StateHistory struct {
	Undo []Transition
	Redo []Transition
}