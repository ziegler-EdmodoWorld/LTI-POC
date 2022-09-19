// Package parser implements a parser for SQL statements
//
// Inspired by go/parser
package parser

// Instruction define a valid SQL statement
type Instruction struct {
	Decls []*Decl
}

// PrettyPrint prints instruction's declarations on console with indentation
func (i Instruction) PrettyPrint() {
	for _, d := range i.Decls {
		d.Stringy(0)
	}
}
