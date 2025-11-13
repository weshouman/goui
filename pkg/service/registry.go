package service

import (
	"sync"

	"github.com/ourorg/goui/pkg/domain"
)

// Registry interface for all registries (following TODO19 core requirements)
type Registry interface {
	Index() interface{}
	Add(...interface{})
	Invalidate()
}

// Concrete registries

type StateRegistry struct {
	mu     sync.Mutex
	states []domain.State
}

func NewStateRegistry() *StateRegistry {
	return &StateRegistry{}
}

func (r *StateRegistry) Index() map[int]domain.State {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[int]domain.State, len(r.states))
	for _, s := range r.states {
		out[s.ID] = s
	}
	return out
}

func (r *StateRegistry) Add(states ...domain.State) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.states = append(r.states, states...)
}

func (r *StateRegistry) Invalidate() {
	r.mu.Lock()
	r.states = nil
	r.mu.Unlock()
}

// Backward compatible helper
func (r *StateRegistry) GetStates() []domain.State {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.states
}

type ModeRegistry struct {
	mu    sync.Mutex
	modes []domain.Mode
}

func NewModeRegistry() *ModeRegistry {
	return &ModeRegistry{}
}

func (r *ModeRegistry) Index() map[int]domain.Mode {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[int]domain.Mode, len(r.modes))
	for _, m := range r.modes {
		out[m.ID] = m
	}
	return out
}

func (r *ModeRegistry) Add(modes ...domain.Mode) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modes = append(r.modes, modes...)
}

func (r *ModeRegistry) Invalidate() {
	r.mu.Lock()
	r.modes = nil
	r.mu.Unlock()
}

// Backward compatible helper
func (r *ModeRegistry) GetModes() []domain.Mode {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.modes
}

type CommandRegistry struct {
	mu       sync.Mutex
	commands []*domain.Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{}
}

func (r *CommandRegistry) Index() map[string]*domain.Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]*domain.Command, len(r.commands))
	for _, c := range r.commands {
		for _, a := range c.Aliases {
			out[a] = c
		}
	}
	return out
}

func (r *CommandRegistry) Add(cmds ...*domain.Command) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands = append(r.commands, cmds...)
}

func (r *CommandRegistry) Invalidate() {
	r.mu.Lock()
	r.commands = nil
	r.mu.Unlock()
}

// Backward compatible helper
func (r *CommandRegistry) GetCommands() []*domain.Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.commands
}

// Backward compatible facade used by apps today
type RegistryFacade struct {
	states   *StateRegistry
	modes    *ModeRegistry
	commands *CommandRegistry
}

func NewRegistry() *RegistryFacade {
	return &RegistryFacade{
		states:   NewStateRegistry(),
		modes:    NewModeRegistry(),
		commands: NewCommandRegistry(),
	}
}

func (r *RegistryFacade) AddStates(states ...domain.State) {
	r.states.Add(states...)
}

func (r *RegistryFacade) AddModes(modes ...domain.Mode) {
	r.modes.Add(modes...)
}

func (r *RegistryFacade) AddCommands(cmds ...*domain.Command) {
	r.commands.Add(cmds...)
}

func (r *RegistryFacade) GetStates() []domain.State {
	return r.states.GetStates()
}

func (r *RegistryFacade) GetModes() []domain.Mode {
	return r.modes.GetModes()
}

func (r *RegistryFacade) GetCommands() []*domain.Command {
	return r.commands.GetCommands()
}

func (r *RegistryFacade) InvalidateStates() {
	r.states.Invalidate()
}

func (r *RegistryFacade) InvalidateModes() {
	r.modes.Invalidate()
}

func (r *RegistryFacade) InvalidateCommands() {
	r.commands.Invalidate()
}

// Helpers for Engine wiring
func (r *RegistryFacade) StateRegistry() *StateRegistry {
	return r.states
}

func (r *RegistryFacade) ModeRegistry() *ModeRegistry {
	return r.modes
}

func (r *RegistryFacade) CommandRegistry() *CommandRegistry {
	return r.commands
}