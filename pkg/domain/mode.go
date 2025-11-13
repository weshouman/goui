package domain

// Mode represents application input modes
type Mode struct {
	ID   int
	Name string
}

// Mode constants
const (
	ModeNormal = iota
	ModeSearch
	ModeCommand
)