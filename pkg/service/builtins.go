package service

import (
	"sort"
	"strings"

	"github.com/ourorg/goui/pkg/domain"
	"github.com/ourorg/goui/pkg/spec"
)

const (
	stateAliases = -101
	stateHelp    = -102
)

func RegisterBuiltins(reg *RegistryFacade, quit func(), showHelp func(), showAliases func()) {
	// Add built-in states
	reg.AddStates(
		domain.State{
			ID:            stateAliases,
			ShortNameTmpl: "Aliases",
			LayoutKind:    domain.DisplayTable,
			Args: map[string]interface{}{
				"title":   "Command Aliases",
				"headers": []string{"Aliases", "Shortcuts", "Template", "From", "To"},
			},
		},
		domain.State{
			ID:            stateHelp,
			ShortNameTmpl: "Help",
			LayoutKind:    domain.DisplayText,
			Args:          map[string]interface{}{},
		},
	)

	reg.AddCommands(
		&domain.Command{
			Aliases:    []string{"q", "quit"},
			FromStates: []int{domain.StateAny},
			ToStates:   []int{domain.StateSame},
			Handler: func(ctx *domain.Ctx, _ []string) (string, error) {
				if quit != nil { quit() }
				return "Quitting...", nil
			},
		},
		&domain.Command{
			Aliases:    []string{"?", "h", "help"},
			FromStates: []int{domain.StateAny},
			ToStates:   []int{domain.StateSame},
			Handler: func(ctx *domain.Ctx, _ []string) (string, error) {
				if showHelp != nil { showHelp() }
				return "Help shown", nil
			},
		},
		&domain.Command{
			Aliases:    []string{"a", "aliases"},
			FromStates: []int{domain.StateAny},
			ToStates:   []int{stateAliases},
			Handler: func(ctx *domain.Ctx, _ []string) (string, error) {
				// Build aliases table data directly
				commands := ctx.Registry.GetCommands()
				states := ctx.Registry.GetStates()
				headers, rows := BuildAliasesTableModelWithShortcuts(commands, states)

				// Convert rows to entries for the new architecture
				var entries []spec.Entry
				for _, row := range rows {
					entries = append(entries, spec.Entry{
						ID:     row[0], // Use first column (aliases) as ID
						Values: row,
					})
				}

				ctx.State.SetNextState(stateAliases, func(a map[string]interface{}) {
					a["title"] = "Command Aliases"
					a["headers"] = headers
					a["entries"] = entries
				})
				return "Aliases listed", nil
			},
		},
	)
}

// Helper to build a simple table model from commands
func BuildAliasesTableModel(commands []*domain.Command) (headers []string, rows [][]string) {
	headers = []string{"Aliases", "Template", "From", "To"}
	// deterministic order
	type row struct{ a, t, f, to string }
	list := make([]row, 0, len(commands))
	for _, c := range commands {
		a := strings.Join(c.Aliases, ", ")
		f := "any"
		if len(c.FromStates) > 0 {
			parts := make([]string, len(c.FromStates))
			for i, s := range c.FromStates { parts[i] = stateName(s) }
			f = strings.Join(parts, "|")
		}
		to := ""
		if len(c.ToStates) > 0 {
			parts := make([]string, len(c.ToStates))
			for i, s := range c.ToStates { parts[i] = stateName(s) }
			to = strings.Join(parts, "|")
		}
		list = append(list, row{a: a, t: c.CmdTmpl, f: f, to: to})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].a < list[j].a })
	for _, r := range list {
		rows = append(rows, []string{r.a, r.t, r.f, r.to})
	}
	return
}

func stateName(id int) string {
	if id == domain.StateAny { return "any" }
	if id == domain.StateSame { return "same" }
	return string(rune('0' + id)) // cheap label; apps know their own mapping
}

// NEW: same model, plus a Shortcuts column computed from state bulk_keys
func BuildAliasesTableModelWithShortcuts(
	commands []*domain.Command,
	states []domain.State,
) (headers []string, rows [][]string) {

	headers = []string{"Aliases", "Shortcuts", "Template", "From", "To"}

	// alias -> set(keys)
	aliasToKeys := map[string]map[string]bool{}
	for _, st := range states {
		if st.Args == nil {
			continue
		}
		if m, ok := st.Args["bulk_keys"].(map[string]string); ok && m != nil {
			for key, alias := range m {
				if alias == "" {
					continue
				}
				set := aliasToKeys[alias]
				if set == nil {
					set = map[string]bool{}
					aliasToKeys[alias] = set
				}
				set[key] = true
			}
		}
	}

	type row struct{ a, k, t, f, to string }
	var list []row

	for _, c := range commands {
		// aliases text
		a := strings.Join(c.Aliases, ", ")

		// union of keys for any alias of this command
		keysSet := map[string]bool{}
		for _, al := range c.Aliases {
			for k := range aliasToKeys[al] {
				keysSet[k] = true
			}
		}
		var keys []string
		for k := range keysSet {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		k := strings.Join(keys, ", ")

		// from states
		f := "any"
		if len(c.FromStates) > 0 {
			parts := make([]string, len(c.FromStates))
			for i, s := range c.FromStates {
				if s == domain.StateAny {
					parts[i] = "any"
				} else if s == domain.StateSame {
					parts[i] = "same"
				} else {
					// simple numeric id label kept like before
					parts[i] = string(rune('0' + s))
				}
			}
			f = strings.Join(parts, "|")
		}

		// to states
		to := ""
		if len(c.ToStates) > 0 {
			parts := make([]string, len(c.ToStates))
			for i, s := range c.ToStates {
				if s == domain.StateAny {
					parts[i] = "any"
				} else if s == domain.StateSame {
					parts[i] = "same"
				} else {
					parts[i] = string(rune('0' + s))
				}
			}
			to = strings.Join(parts, "|")
		}

		list = append(list, row{a: a, k: k, t: c.CmdTmpl, f: f, to: to})
	}

	sort.Slice(list, func(i, j int) bool { return list[i].a < list[j].a })
	for _, r := range list {
		rows = append(rows, []string{r.a, r.k, r.t, r.f, r.to})
	}
	return
}