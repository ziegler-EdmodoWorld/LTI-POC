package engine

// Table is defined by a name and attributes
// A table with data is called a Relation
type Table struct {
	name       string
	attributes []Attribute
}

// NewTable initializes a new Table
func NewTable(name string) *Table {
	t := &Table{
		name: name,
	}

	return t
}

// AddAttribute is used by CREATE TABLE and ALTER TABLE
// Want to check that name isn't already taken
func (t *Table) AddAttribute(attr Attribute) error {
	t.attributes = append(t.attributes, attr)
	return nil
}

// String returns a printable string with table name and attributes
func (t Table) String() string {
	stringy := t.name + " ("
	for i, a := range t.attributes {
		if i != 0 {
			stringy += " | "
		}
		stringy += a.name + " " + a.typeName
	}
	stringy += ")"
	return stringy
}
