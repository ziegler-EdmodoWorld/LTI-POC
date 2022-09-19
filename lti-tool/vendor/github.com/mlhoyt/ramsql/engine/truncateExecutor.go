package engine

import (
	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func truncateExecutor(e *Engine, trDecl *parser.Decl, conn protocol.EngineConn) error {
	log.Debug("truncateExecutor")

	// get tables to be deleted
	table := NewTable(trDecl.Decl[0].Lexeme)

	return truncateTable(e, table, conn)
}
