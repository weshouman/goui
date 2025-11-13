package execx

import (
	"bytes"
	"text/template"
	"time"
)

type demoExec struct{ cfg Config }

func NewDemo(cfg Config) Executor {
	if cfg.Timeout == 0 { cfg.Timeout = 10 * time.Second }
	return &demoExec{cfg: cfg}
}

func (e *demoExec) Mode() Mode { return ModeDemo }

func (e *demoExec) Run(argv ...string) (Result, error) {
	if e.cfg.DemoLatency > 0 { time.Sleep(e.cfg.DemoLatency) }
	return Result{
		Stdout:   "demo: " + join(argv),
		Stderr:   "",
		ExitCode: 0,
	}, nil
}

func (e *demoExec) RunTemplate(tmpl string, data map[string]interface{}) (Result, error) {
	t, err := template.New("cmd").Parse(tmpl)
	if err != nil { return Result{}, err }
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil { return Result{}, err }
	return e.Run(buf.String())
}

func join(a []string) string {
	if len(a) == 0 { return "" }
	s := a[0]
	for i := 1; i < len(a); i++ { s += " " + a[i] }
	return s
}