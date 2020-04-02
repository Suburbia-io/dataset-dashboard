package dbgen

// ----------------------------------------------------------------------------

type rawCol struct {
	DBName         string // The column's database name.
	GoType         string // The Go type.
	GoName         string // The Go name of the column.
	JSONName       string // JSON name for the column (struct tag).
	PK             bool   // This is a primary key field.
	AutoUUID       bool   // Auto generated UUID if empty string.
	Sanitizer      string // From the sanitize helper package.
	Validator      string // From the validate helper package.
	ViewOnly       bool   // This field isn't updated.
	NoUpdate       bool   // This field isn't modified by the update function.
	NoDirectUpdate bool   // This field doesn't have a function for direct update.
}

func (r *rawCol) Insertable() bool {
	return !r.ViewOnly
}

func (r *rawCol) Updatable() bool {
	return !(r.PK || r.ViewOnly || r.NoUpdate)
}

func (r *rawCol) DirectUpdatable() bool {
	return !(r.PK || r.ViewOnly || r.NoDirectUpdate)
}

// IsPointer returns true if the column's GoType is a pointer.
func (c *rawCol) IsPointer() bool {
	return c.GoType[0] == '*'
}

// ----------------------------------------------------------------------------

type Col struct {
	c *rawCol
}

func newCol(dbName, goType string) Col {
	return Col{
		c: &rawCol{
			DBName:         dbName,
			GoType:         goType,
			GoName:         SQLNameToGo(dbName),
			JSONName:       SQLNameToJSON(dbName),
			NoDirectUpdate: true,
		},
	}
}

func (c Col) GoName(s string) Col {
	c.c.GoName = s
	return c
}

func (c Col) JSONName(s string) Col {
	c.c.JSONName = s
	return c
}

func (c Col) Validator(s string) Col {
	c.c.Validator = s
	return c
}

func (c Col) Sanitizer(s string) Col {
	c.c.Sanitizer = s
	return c
}

func (c Col) PK() Col {
	c.c.PK = true
	return c
}

func (c Col) AutoUUID() Col {
	c.c.AutoUUID = true
	return c
}

func (c Col) ViewOnly() Col {
	c.c.ViewOnly = true
	c.c.NoUpdate = true
	c.c.NoDirectUpdate = true
	return c
}

func (c Col) NoUpdate() Col {
	c.c.NoUpdate = true
	return c
}

func (c Col) DirectUpdate() Col {
	c.c.NoDirectUpdate = false
	return c
}
