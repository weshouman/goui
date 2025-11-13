package domain

// StateWriter is the write-only interface for state operations
// This is duplicated from service package to avoid import cycles
type StateWriter interface {
	SetNextState(id int, mutateArgs func(map[string]interface{})) error
	Push(id int) error
	Pop() error
}