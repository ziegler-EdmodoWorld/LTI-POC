package engine

import (
	"fmt"

	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func dropExecutor(e *Engine, dropDecl *parser.Decl, conn protocol.EngineConn) error {
	// Action Parameters
	var allowNotFound = false
	var tableName string

	// Process Decls

	// Required: TABLE
	if dropDecl.Decl == nil ||
		len(dropDecl.Decl) != 1 ||
		dropDecl.Decl[0].Token != parser.TableToken {
		return fmt.Errorf("unexpected drop arguments")
	}

	// Optional: IF EXISTS
	tableNameTokenIndex := 0
	if dropDecl.Decl[0].Decl[0].Token == parser.IfToken {
		allowNotFound = true
		tableNameTokenIndex = 1
	}

	// Required: <TABLE-NAME>
	tableName = dropDecl.Decl[0].Decl[tableNameTokenIndex].Lexeme

	// Pre-Action/s
	r := e.relation(tableName)
	if r == nil {
		if allowNotFound {
			return conn.WriteResult(0, 1)
		}

		return fmt.Errorf("relation '%s' not found", tableName)
	}

	// Action/s
	e.drop(tableName)

	// Post-Action/s
	// None

	return conn.WriteResult(0, 1)
}
