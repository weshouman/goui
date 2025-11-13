package service

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ourorg/goui/pkg/domain"
)

type CmdHistoryEntry struct {
	Cmd   string
	Count int
	Last  time.Time
}

type CmdHistory struct {
	entries map[string]*CmdHistoryEntry
	mu      sync.Mutex
}

func newCmdHistory() *CmdHistory {
	return &CmdHistory{entries: map[string]*CmdHistoryEntry{}}
}

func (h *CmdHistory) Touch(cmd string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := h.entries[cmd]
	if e == nil {
		e = &CmdHistoryEntry{Cmd: cmd}
		h.entries[cmd] = e
	}
	e.Count++
	e.Last = time.Now()
}

func (h *CmdHistory) Entries() map[string]*CmdHistoryEntry {
	return h.entries
}

type CommandService struct {
	cmdReg *CommandRegistry
	hist   *CmdHistory
	// Engine wiring
	ctxBuilder func() *domain.Ctx
}

func NewCommandService(reg *CommandRegistry, ctxBuilder func() *domain.Ctx) *CommandService {
	return &CommandService{
		cmdReg:     reg,
		hist:       newCmdHistory(),
		ctxBuilder: ctxBuilder,
	}
}

func (s *CommandService) TouchHistory(cmd string) {
	s.hist.Touch(cmd)
}

func (s *CommandService) Suggestions(prefix string) []string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return nil // prevent suggestions on empty input
	}
	seen := map[string]bool{}
	var out []string

	// history first
	type row struct {
		cmd  string
		last time.Time
	}
	var hs []row
	for _, e := range s.hist.Entries() {
		if strings.HasPrefix(e.Cmd, prefix) {
			hs = append(hs, row{e.Cmd, e.Last})
		}
	}
	sort.Slice(hs, func(i, j int) bool {
		return hs[i].last.After(hs[j].last)
	})
	for _, r := range hs {
		out = append(out, r.cmd)
		seen[r.cmd] = true
	}

	// then all aliases
	var rest []string
	for _, c := range s.cmdReg.Index() {
		for _, a := range c.Aliases {
			if strings.HasPrefix(a, prefix) && !seen[a] {
				rest = append(rest, a)
			}
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}

func (s *CommandService) Autocomplete(prefix string) []string {
	var res []string
	for _, c := range s.cmdReg.Index() {
		for _, a := range c.Aliases {
			if len(prefix) == 0 || strings.HasPrefix(a, prefix) {
				res = append(res, a)
			}
		}
	}
	return res
}

func (s *CommandService) Resolve(alias string) (*domain.Command, bool) {
	for _, c := range s.cmdReg.Index() {
		for _, a := range c.Aliases {
			if a == alias {
				return c, true
			}
		}
	}
	return nil, false
}

func (s *CommandService) Dispatch(alias string, args []string) (string, error) {
	cmd, ok := s.Resolve(alias)
	if !ok {
		return "Unknown command: " + alias, nil
	}
	s.TouchHistory(alias)

	ctx := s.ctxBuilder()
	if cmd.Handler == nil {
		return "Executing mock: " + cmd.CmdTmpl, nil
	}
	return cmd.Handler(ctx, args)
}