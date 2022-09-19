package engine

import (
	"errors"

	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func createExecutor(e *Engine, createDecl *parser.Decl, conn protocol.EngineConn) error {
	if len(createDecl.Decl) == 0 {
		return errors.New("Parsing failed, no declaration after CREATE")
	}

	if _, ok := e.opsExecutors[createDecl.Decl[0].Token]; !ok {
		return errors.New("Parsing failed, after CREATE unkown token " + createDecl.Decl[0].Lexeme)
	}

	return e.opsExecutors[createDecl.Decl[0].Token](e, createDecl.Decl[0], conn)
}
