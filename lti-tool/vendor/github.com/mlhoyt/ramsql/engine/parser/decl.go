// Package parser implements a parser for SQL statements
//
// Inspired by go/parser
package parser

import (
	"fmt"

	"github.com/mlhoyt/ramsql/engine/log"
)

// Decl structure is the node to statement declaration tree
type Decl struct {
	Token  int
	Lexeme string
	Decl   []*Decl
}

// NewDecl initialize a Decl struct from a given token
func NewDecl(t Token) *Decl {
	return &Decl{
		Token:  t.Token,
		Lexeme: t.Lexeme,
	}
}

// Add creates a new leaf with given Decl
func (d *Decl) Add(subDecl *Decl) {
	d.Decl = append(d.Decl, subDecl)
}

// Stringy prints the declaration tree in console
func (d Decl) Stringy(depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent = fmt.Sprintf("%s    ", indent)
	}

	log.Debug("%s|-> %s\n", indent, d.Lexeme)
	for _, subD := range d.Decl {
		subD.Stringy(depth + 1)
	}
}
