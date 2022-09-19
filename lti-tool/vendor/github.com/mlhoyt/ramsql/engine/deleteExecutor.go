package engine

import (
	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func deleteExecutor(e *Engine, deleteDecl *parser.Decl, conn protocol.EngineConn) error {
	log.Debug("deleteExecutor")

	// get tables to be deleted
	tables := fromExecutor(deleteDecl.Decl[0])

	// If len is 1, it means no predicates so truncate table
	if len(deleteDecl.Decl) == 1 {
		return truncateTable(e, tables[0], conn)
	}

	// get WHERE declaration
	predicates, err := whereExecutor(deleteDecl.Decl[1], tables[0].name)
	if err != nil {
		return err
	}

	// and delete
	return deleteRows(e, tables, conn, predicates)
}
