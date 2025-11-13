package domain

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/ourorg/goui/pkg/execx"
)

// Command represents a command that can be executed in the application
type Command struct {
	Aliases []string
	CmdTmpl string
	Args    []string

	FromStates []int
	ToStates   []int
	NextStateLogic func(currState int) int

	FromMode int
	ToMode   int
	ModeTransitionLogic func(currentMode int) int

	// New: project-defined mock or real handler
	Handler func(ctx *Ctx, args []string) (string, error)

	// Framework-provided info sink
	SetInfo func(string)
}

// Ctx provides context for command execution
type Ctx struct {
	CurrentStateID int
	Registry       RegistryReader

	// Execution support (preserved)
	Exec     execx.Executor
	ExecMode execx.Mode

	// TODO19 architecture: StateWriter for clean state mutations
	State StateWriter

	// optional future hooks
	// SetExecMode func(mode execx.Mode, cfg execx.Config) error
}

// RegistryReader provides read-only access to registry data
type RegistryReader interface {
	GetStates() []State
	GetCommands() []*Command
}

// IsAvailable checks if a command is available in a given state.
func (c *Command) IsAvailable(stateID int) bool {
	if len(c.FromStates) == 0 {
		logrus.Debugf("Command not available - no FromStates defined")
		return false
	}
	for _, sID := range c.FromStates {
		if sID == -1 || sID == stateID {
			logrus.Debugf("Command available in state %d", stateID)
			return true
		}
	}
	logrus.Debugf("Command not available in state %d", stateID)
	return false
}

func (c *Command) NextState(currState int) int {
	if c.NextStateLogic != nil {
		logrus.Debugf("Command has a NextStateLogic")
		nextState := c.NextStateLogic(currState)
		logrus.Debugf("NextStateLogic returned state: %d", nextState)
		return nextState
	}
	if len(c.ToStates) > 0 {
		logrus.Debugf("Command has %d ToStates", len(c.ToStates))
		logrus.Debugf("Using first ToState: %d", c.ToStates[0])
		return c.ToStates[0] // Return the first state by default.
	}
	logrus.Debugf("Command does not change the state")
	logrus.Debugf("Staying in current state: %d", currState)
	return currState // Stay in the current state if no ToStates are defined, also StateSame would fall here
}

func FindCommandByAlias(alias string, commands []*Command) *Command {
	logrus.Debugf("Searching for command with alias: %s", alias)
	for _, cmd := range commands {
		for _, a := range cmd.Aliases {
			if a == alias {
				logrus.Debugf("Found command with template: %s", cmd.CmdTmpl)
				return cmd
			}
		}
	}
	logrus.Debugf("No command found for alias: %s", alias)
	return nil
}

func (c *Command) Execute() {
	// Here, execute the command using CmdTmpl and Args.
	logrus.Debugf("Executing command template: %s", c.CmdTmpl)
	logrus.Debugf("Command arguments: %v", c.Args)
	// TODO: expand the command template and run it
	// Support remote execution
	// Correctly handle special commands that are relevant
	//   to the app itself like searching state saving
	formattedMessage := fmt.Sprintf("Executing: %s with args: %v\n", c.CmdTmpl, c.Args)
	logrus.Info(formattedMessage)
	c.SetInfo(formattedMessage)
}

// ParseInput splits command input into alias and args
func ParseInput(text string) (alias string, args []string) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}