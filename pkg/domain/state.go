package domain

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ourorg/goui/pkg/util"
)

// State represents an application state with its display configuration
type State struct {
	ID            int
	Name          string // Human-readable name of the state
	ShortNameTmpl string // The name could have the args used
	LongNameTmpl  string

	// StateDisplay is reserved in case further common info for the state display
	// is to be considered, it could be simplified though to StateLayoutKind
	// if only the type is to be used
	LayoutKind int

	Args map[string]interface{}

	// New TODO19: Caching support as per core requirements
	CacheTTLSeconds int       // 0 for no caching, N seconds for cache validity
	ComputedAt      time.Time // When the state data was last computed

	// Optional selection handler for later
	OnSelect func(item string)
}

// Display layout constants
const (
	DisplayText int = iota
	DisplayTable
	DisplayList
	DisplayTree
	DisplayTreeText
	DisplayTreeTable
	DisplayListText
	DisplayListTable
)

// State constants
const (
	StateAny  = -1
	StateSame = -2
)

var displayTypeToString = map[int]string{
	DisplayText:      "DisplayText",
	DisplayTable:     "DisplayTable",
	DisplayList:      "DisplayList",
	DisplayTree:      "DisplayTree",
	DisplayTreeText:  "DisplayTreeText",
	DisplayTreeTable: "DisplayTreeTable",
	DisplayListText:  "DisplayListText",
	DisplayListTable: "DisplayListTable",
}

func LayoutKindToString(t int) string {
	if str, ok := displayTypeToString[t]; ok {
		return str
	}
	return "Unknown"
}

func (s *State) ShortName() string {
	name := util.ProcessTemplate(s.ShortNameTmpl, s.Args)
	logrus.Debugf("Generated short name for state %d: %s", s.ID, name)
	return name
}

func (s *State) LongName() string {
	if s.LongNameTmpl == "" {
		s.LongNameTmpl = s.ShortNameTmpl
	}
	name := util.ProcessTemplate(s.LongNameTmpl, s.Args)
	logrus.Debugf("Generated long name for state %d: %s", s.ID, name)
	return name
}

func GetStateByID(states []State, id int) (*State, error) {
	for _, state := range states {
		if state.ID == id {
			return &state, nil
		}
	}
	return nil, fmt.Errorf("no state found with ID: %d", id)
}