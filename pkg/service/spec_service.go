package service

import (
	"strings"

	"github.com/ourorg/goui/pkg/domain"
	"github.com/ourorg/goui/pkg/spec"
)

type SpecService struct{}

func NewSpecService() *SpecService { return &SpecService{} }

// helpers kept from old StateManager
func containsCI(hay, needle string) bool {
	return strings.Contains(strings.ToLower(hay), strings.ToLower(needle))
}

func filterRows(rows [][]string, term string) [][]string {
	if term == "" { return rows }
	var out [][]string
	for _, r := range rows {
		for _, c := range r {
			if containsCI(c, term) { out = append(out, r); break }
		}
	}
	return out
}

func filterEntries(entries []spec.Entry, term string) []spec.Entry {
	if term == "" { return entries }
	var out []spec.Entry
	for _, e := range entries {
		for _, v := range e.Values {
			if containsCI(v, term) { out = append(out, e); break }
		}
	}
	return out
}

func filterList(items []spec.ListItem, term string) []spec.ListItem {
	if term == "" { return items }
	var out []spec.ListItem
	for _, it := range items {
		if containsCI(it.Main, term) || containsCI(it.Secondary, term) {
			out = append(out, it)
		}
	}
	return out
}

func filterText(body, term string) string {
	if term == "" { return body }
	var out []string
	for _, line := range strings.Split(body, "\n") {
		if containsCI(line, term) { out = append(out, line) }
	}
	if len(out) == 0 { return "No matches" }
	return strings.Join(out, "\n")
}

func (s *SpecService) BuildSpec(st *domain.State) spec.Spec {
	if st == nil {
		return spec.Spec{Kind: spec.KindText, Text: &spec.Text{Body: "no state"}}
	}
	searchTerm, _ := st.Args["searchTerm"].(string)

	switch st.LayoutKind {
	case domain.DisplayTable:
		title := "Items"
		if t, ok := st.Args["title"].(string); ok && t != "" { title = t }

		headers := []string{"Name"}
		if h, ok := st.Args["headers"].([]string); ok && len(h) > 0 { headers = h }

		// Build entries, keeping Rows as a fallback for older screens
		var entries []spec.Entry
		if en, ok := st.Args["entries"].([]spec.Entry); ok && len(en) > 0 {
			entries = en
		} else {
			// Make entries from rows, derive IDs from id_col if present
			idCol := 0
			if v, ok := st.Args["id_col"].(int); ok { idCol = v }
			rows := [][]string{{"item_1"}, {"item_2"}}
			if r, ok := st.Args["rows"].([][]string); ok && len(r) > 0 { rows = r }

			// Apply filtering first on rows if we use them
			rows = filterRows(rows, searchTerm)

			for _, r := range rows {
				id := ""
				if idCol >= 0 && idCol < len(r) { id = r[idCol] }
				entries = append(entries, spec.Entry{ID: id, Values: r})
			}
		}

		// Final filter on entries to support both paths
		entries = filterEntries(entries, searchTerm)

		// Selection comes from state args, optional
		var sel []string
		if v, ok := st.Args["selection"].([]string); ok && len(v) > 0 {
			sel = append(sel, v...)
		}

		return spec.Spec{
			Kind: spec.KindTable,
			Table: &spec.Table{
				Title:    title,
				Headers:  headers,
				Entries:  entries,
				// Keep Rows for compatibility so older renderers still show something
				Rows:     valuesFromEntries(entries),
			},
			Selection: sel,
		}

	case domain.DisplayList:
		items := []spec.ListItem{
			{Main: "Default 1", Secondary: "Description 1"},
			{Main: "Default 2", Secondary: "Description 2"},
		}
		if li, ok := st.Args["list"].([]spec.ListItem); ok && len(li) > 0 {
			items = li
		}

		items = filterList(items, searchTerm)

		// Optional list selection, we use item Main as ID by default
		var sel []string
		if v, ok := st.Args["selection"].([]string); ok && len(v) > 0 {
			sel = append(sel, v...)
		}

		return spec.Spec{
			Kind: spec.KindList,
			List: &spec.List{
				Title: "List",
				Items: items,
			},
			Selection: sel,
		}

	default:
		body := "Default text content\nApps should provide custom content"
		if s, ok := st.Args["text"].(string); ok && s != "" { body = s }
		body = filterText(body, searchTerm)
		return spec.Spec{
			Kind: spec.KindText,
			Text: &spec.Text{Title: st.ShortName(), Body: body},
		}
	}
}

func (s *SpecService) ApplyFilter(sp spec.Spec, args map[string]interface{}) spec.Spec {
	// optional custom filtering stage, keeping it as passthrough for now
	return sp
}

func valuesFromEntries(es []spec.Entry) [][]string {
	out := make([][]string, 0, len(es))
	for _, e := range es { out = append(out, e.Values) }
	return out
}