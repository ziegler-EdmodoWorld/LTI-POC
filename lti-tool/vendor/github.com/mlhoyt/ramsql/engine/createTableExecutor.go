package engine

import (
	"fmt"

	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func createTableExecutor(e *Engine, tableDecl *parser.Decl, conn protocol.EngineConn) error {
	if len(tableDecl.Decl) == 0 {
		return fmt.Errorf("parsing failed, malformed CREATE TABLE query")
	}

	// Fetch constrainit (i.e: "IF EXISTS")
	i := 0
	for i < len(tableDecl.Decl) {
		if e.opsExecutors[tableDecl.Decl[i].Token] != nil {
			if err := e.opsExecutors[tableDecl.Decl[i].Token](e, tableDecl.Decl[i], conn); err != nil {
				return err
			}
		} else {
			break
		}

		i++
	}

	// Check if table does not exists
	r := e.relation(tableDecl.Decl[i].Lexeme)
	if r != nil {
		return fmt.Errorf("table %s already exists", tableDecl.Decl[i].Lexeme)
	}

	// Fetch table name
	t := NewTable(tableDecl.Decl[i].Lexeme)

	// Fetch attributes
	i++
	for i < len(tableDecl.Decl) {
		attr, err := parseAttribute(tableDecl.Decl[i])
		if err != nil {
			return err
		}
		err = t.AddAttribute(attr)
		if err != nil {
			return err
		}
		i++
	}

	e.relations[t.name] = NewRelation(t)
	conn.WriteResult(0, 1)
	return nil
}
