package execx

import (
	"time"
)

type Mode int

const (
	ModeDemo Mode = iota
	ModeLocal
	ModeSSH
)

func (m Mode) String() string {
	switch m {
	case ModeDemo:
		return "demo"
	case ModeLocal:
		return "local"
	case ModeSSH:
		return "ssh"
	default:
		return "unknown"
	}
}

type Config struct {
	Mode       Mode
	Timeout    time.Duration // default 10s
	// ssh settings for ModeSSH
	SSHHost    string   // host or host:port
	SSHUser    string   // optional, empty uses default
	SSHOptions []string // extra ssh -o options
	// demo behavior
	DemoLatency time.Duration
}

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Executor is the unified command runner
type Executor interface {
	Mode() Mode
	Run(argv ...string) (Result, error)
	RunTemplate(tmpl string, data map[string]interface{}) (Result, error)
}