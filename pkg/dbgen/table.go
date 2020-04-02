package dbgen

type GetBy struct {
	Name string
	Cols []*rawCol
}

type rawTable struct {
	Package       string
	Table         string
	View          string
	RowType       string
	Cols          []*rawCol
	GetBys        []GetBy
	NoUpdateTests bool // Don't generate tests for update functions.

	colMap  map[string]*rawCol
	counter int
}

func (t *rawTable) HasValidator() bool {
	for _, col := range t.Cols {
		if col.Validator != "" {
			return true
		}
	}
	return false
}

func (t *rawTable) HasSanitizer() bool {
	for _, col := range t.Cols {
		if col.Sanitizer != "" {
			return true
		}
	}
	return false
}

func (t *rawTable) HasAutoUUID() bool {
	for _, col := range t.Cols {
		if col.AutoUUID {
			return true
		}
	}
	return false
}

func (t *rawTable) Updatable() bool {
	for _, col := range t.Cols {
		if col.Updatable() {
			return true
		}
	}
	return false
}

func (t *rawTable) InsertCols() (l []*rawCol) {
	for _, c := range t.Cols {
		if c.Insertable() {
			l = append(l, c)
		}
	}
	return l
}

func (t *rawTable) KeyCols() (l []*rawCol) {
	for _, c := range t.Cols {
		if c.PK {
			l = append(l, c)
		}
	}
	return l
}

func (t *rawTable) UpdateCols() (l []*rawCol) {
	for _, c := range t.Cols {
		if c.Updatable() {
			l = append(l, c)
		}
	}
	return l
}

func (t *rawTable) DirectUpdateCols() (l []*rawCol) {
	for _, c := range t.Cols {
		if c.DirectUpdatable() {
			l = append(l, c)
		}
	}
	return l
}

func (t *rawTable) ResetCounter() string {
	t.counter = 0
	return ""
}

func (t *rawTable) Counter() int {
	t.counter++
	return t.counter
}

// ----------------------------------------------------------------------------

type Table struct {
	t *rawTable
}

func NewTable(name string) Table {
	return Table{
		t: &rawTable{
			Table:   name,
			View:    name,
			RowType: TableNameToRowType(name),
			Cols:    []*rawCol{},
			GetBys:  []GetBy{},
			colMap:  map[string]*rawCol{},
		},
	}
}

func (t Table) View(s string) Table {
	t.t.View = s
	return t
}

func (t Table) RowType(s string) Table {
	t.t.RowType = s
	return t
}

func (t Table) Col(dbName, goType string) Col {
	c := newCol(dbName, goType)
	t.t.Cols = append(t.t.Cols, c.c)
	t.t.colMap[dbName] = c.c
	return c
}

func (t Table) NoUpdateTests() Table {
	t.t.NoUpdateTests = true
	return t
}

func (t Table) GetBy(name string, cols ...string) Table {
	keyCols := []*rawCol{}
	for _, colName := range cols {
		keyCols = append(keyCols, t.t.colMap[colName])
	}
	t.t.GetBys = append(t.t.GetBys, GetBy{
		Name: name,
		Cols: keyCols,
	})
	return t
}
