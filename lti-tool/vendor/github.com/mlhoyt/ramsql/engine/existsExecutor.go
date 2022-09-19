package engine

import (
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func existsExecutor(e *Engine, tableDecl *parser.Decl, conn protocol.EngineConn) error {
	return nil
}
