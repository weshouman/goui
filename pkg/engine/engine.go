package engine

import (
	"errors"

	"github.com/ourorg/goui/pkg/domain"
	"github.com/ourorg/goui/pkg/execx"
	"github.com/ourorg/goui/pkg/service"
	"github.com/ourorg/goui/pkg/spec"
)

type Options struct {
	Info       func(string)
	ExecMode   execx.Mode
	ExecConfig execx.Config
}

type Engine struct {
	// registries
	stateReg   *service.StateRegistry
	modeReg    *service.ModeRegistry
	cmdReg     *service.CommandRegistry

	// providers
	specService    service.SpecProvider
	stateService   service.StateProvider
	modeService    service.ModeProvider
	commandService service.CommandProvider

	// execution
	execMode execx.Mode
	executor execx.Executor

	// info sink
	info func(string)
}

func New(
	sr *service.StateRegistry,
	mr *service.ModeRegistry,
	cr *service.CommandRegistry,
	sp service.SpecProvider,
	st service.StateProvider,
	md service.ModeProvider,
	cp service.CommandProvider,
	opts Options,
) *Engine {
	e := &Engine{
		stateReg:       sr,
		modeReg:        mr,
		cmdReg:         cr,
		specService:    sp,
		stateService:   st,
		modeService:    md,
		commandService: cp,
		info:           opts.Info,
	}

	// executor
	cfg := opts.ExecConfig
	if cfg.Mode == 0 {
		if opts.ExecMode != 0 {
			cfg.Mode = opts.ExecMode
		} else {
			cfg.Mode = execx.ModeDemo
		}
	}
	switch cfg.Mode {
	case execx.ModeLocal:
		e.executor = execx.NewLocal(cfg)
	case execx.ModeSSH:
		e.executor = execx.NewSSH(cfg)
	default:
		e.executor = execx.NewDemo(cfg)
	}
	e.execMode = cfg.Mode

	// wire SetInfo on commands
	for _, c := range cr.Index() {
		c.SetInfo = func(msg string) {
			if e.info != nil {
				e.info(msg)
			}
		}
	}

	// init state
	_ = e.stateService.Init(firstStateID(sr))
	return e
}

func firstStateID(sr *service.StateRegistry) int {
	idx := sr.Index()
	min := int(^uint(0) >> 1) // max int
	for id := range idx {
		if id < min {
			min = id
		}
	}
	if min == int(^uint(0)>>1) {
		return 0
	}
	return min
}

func (e *Engine) CurrentState() *domain.State {
	return e.stateService.Current()
}

func (e *Engine) SetInfo(fn func(string)) {
	e.info = fn
}

func (e *Engine) BuildSpec() spec.Spec {
	return e.specService.BuildSpec(e.CurrentState())
}

// StateCtrl exposes a writer view for UIs and commands.
func (e *Engine) StateCtrl() domain.StateWriter {
	return stateWriter{e.stateService}
}

type stateWriter struct {
	s service.StateProvider
}

func (w stateWriter) SetNextState(id int, mutateArgs func(map[string]interface{})) error {
	return w.s.SetNextState(id, mutateArgs)
}

func (w stateWriter) Push(id int) error {
	return w.s.Push(id)
}

func (w stateWriter) Pop() error {
	return w.s.Pop()
}

func (e *Engine) Execute(alias string, args []string) (string, spec.Spec, error) {
	if alias == "" {
		return "", e.BuildSpec(), errors.New("empty command")
	}

	// Delegate to CommandService for dispatch
	msg, err := e.commandService.Dispatch(alias, args)
	if msg == "" && err != nil {
		msg = "Error: " + err.Error()
	}

	// Handle mode/state transitions through providers
	if cmd, ok := e.commandService.Resolve(alias); ok && cmd.IsAvailable(e.CurrentState().ID) {
		next := cmd.NextState(e.CurrentState().ID)
		_ = e.stateService.SetNextState(next, nil)
	}

	return msg, e.BuildSpec(), err
}

func (e *Engine) Suggestions(prefix string) []string {
	return e.commandService.Suggestions(prefix)
}

func (e *Engine) Autocomplete(prefix string) []string {
	return e.commandService.Autocomplete(prefix)
}

func (e *Engine) SetMode(m int) {
	e.modeService.SetMode(m)
}

func (e *Engine) CurrentMode() int {
	return e.modeService.CurrentMode()
}

// Expose executor for context building
func (e *Engine) Executor() execx.Executor {
	return e.executor
}

func (e *Engine) ExecMode() execx.Mode {
	return e.execMode
}

// Undo/redo operations for TODO19 architecture
func (e *Engine) Undo() bool {
	if e.stateService.Undo() {
		return true
	}
	return false
}

func (e *Engine) Redo() bool {
	if e.stateService.Redo() {
		return true
	}
	return false
}


// Context builder helper for CommandService wiring.
func NewCtxBuilder(e *Engine, regReader domain.RegistryReader) func() *domain.Ctx {
	return func() *domain.Ctx {
		currState := e.CurrentState()
		stateID := 0
		if currState != nil {
			stateID = currState.ID
		}
		return &domain.Ctx{
			CurrentStateID: stateID,
			Registry:       regReader,
			Exec:           e.executor,
			ExecMode:       e.execMode,
			State:          e.StateCtrl(),
		}
	}
}