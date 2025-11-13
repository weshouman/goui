package execx

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

type localExec struct{ cfg Config }

func NewLocal(cfg Config) Executor {
	if cfg.Timeout == 0 { cfg.Timeout = 10 * time.Second }
	return &localExec{cfg: cfg}
}

func (e *localExec) Mode() Mode { return ModeLocal }

func (e *localExec) Run(argv ...string) (Result, error) {
	if len(argv) == 0 { return Result{}, nil }
	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
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

func (e *localExec) RunTemplate(tmpl string, data map[string]interface{}) (Result, error) {
	t, err := template.New("cmd").Parse(tmpl)
	if err != nil {
		// fallback: treat as shell command
		return e.Run("sh", "-c", tmpl)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return Result{}, err
	}
	return e.Run("sh", "-c", buf.String())
}