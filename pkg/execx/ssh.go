package execx

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

// MVP ssh implementation shells out to system ssh for zero deps.
// Later we can switch to golang.org/x/crypto/ssh.
type sshExec struct{ cfg Config }

func NewSSH(cfg Config) Executor {
	if cfg.Timeout == 0 { cfg.Timeout = 10 * time.Second }
	return &sshExec{cfg: cfg}
}

func (e *sshExec) Mode() Mode { return ModeSSH }

func (e *sshExec) Run(argv ...string) (Result, error) {
	if e.cfg.SSHHost == "" { return Result{Stderr: "ssh host not set"}, nil }

	sshArgs := make([]string, 0, 8)
	sshArgs = append(sshArgs, "-o", "BatchMode=yes")
	for _, o := range e.cfg.SSHOptions {
		sshArgs = append(sshArgs, "-o", o)
	}
	target := e.cfg.SSHHost
	if e.cfg.SSHUser != "" { target = e.cfg.SSHUser + "@" + target }
	sshArgs = append(sshArgs, target)

	// command string to run remotely
	cmdStr := strings.Join(argv, " ")
	sshArgs = append(sshArgs, cmdStr)

	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ssh", sshArgs...)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	runErr := cmd.Run()
	exit := 0
	if cmd.ProcessState != nil {
		exit = cmd.ProcessState.ExitCode()
	}
	return Result{
		Stdout:   strings.TrimRight(out.String(), "\n"),
		Stderr:   strings.TrimRight(errb.String(), "\n"),
		ExitCode: exit,
	}, runErr
}

func (e *sshExec) RunTemplate(tmpl string, data map[string]interface{}) (Result, error) {
	t, err := template.New("cmd").Parse(tmpl)
	if err != nil {
		// simplest path
		return e.Run("sh", "-lc", tmpl)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return Result{}, err
	}
	return e.Run("sh", "-lc", buf.String())
}