package spec

type Kind int

const (
	KindText Kind = iota
	KindTable
	KindList
)

// New in TODO19: Enhanced table model with entries and column metadata
type ColMeta struct {
	Type     string
	Nice     int
	MaxWidth int
	Visible  bool
}

// New in TODO19: Entry with stable ID for selection persistence
type Entry struct {
	ID     string
	Values []string
}

type ListItem struct {
	Main      string
	Secondary string
	Shortcut  rune
}

type Table struct {
	Title    string
	Headers  []string

	// Old representation, kept for backwards compatibility
	Rows     [][]string

	// New preferred representation with stable IDs
	Entries  []Entry

	// Optional column schema
	ColSchema []ColMeta
}

type Text struct {
	Title string
	Body  string
}

type List struct {
	Title string
	Items []ListItem
}

type Spec struct {
	Kind  Kind
	Table *Table
	Text  *Text
	List  *List

	// New: selected IDs (table or list)
	Selection []string
}