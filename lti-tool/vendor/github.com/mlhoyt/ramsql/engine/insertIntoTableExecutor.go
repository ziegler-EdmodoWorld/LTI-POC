package engine

import (
	"fmt"

	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

/*
|-> INSERT
    |-> INTO
        |-> user
            |-> last_name
            |-> first_name
            |-> email
    |-> VALUES
        |-> Roullon
        |-> Pierre
        |-> pierre.roullon@gmail.com
*/
func insertIntoTableExecutor(e *Engine, insertDecl *parser.Decl, conn protocol.EngineConn) error {

	// Get table and concerned attributes and write lock it
	r, attributes, err := getRelation(e, insertDecl.Decl[0])
	if err != nil {
		return err
	}
	r.Lock()
	defer r.Unlock()

	// Check for RETURNING clause
	var returnedID string
	if len(insertDecl.Decl) > 2 {
		for i := range insertDecl.Decl {
			if insertDecl.Decl[i].Token == parser.ReturningToken {
				returnedID = insertDecl.Decl[i].Lexeme
				break
			}
		}
	}

	// Create a new tuple with values
	id, err := insert(r, attributes, insertDecl.Decl[1].Decl, returnedID)
	if err != nil {
		return err
	}

	// if RETURNING decl is not present
	if returnedID != "" {
		conn.WriteRowHeader([]string{returnedID})
		conn.WriteRow([]string{fmt.Sprintf("%v", id)})
		conn.WriteRowEnd()
	} else {
		conn.WriteResult(id, 1)
	}
	return nil
}
