package engine

import (
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

func grantExecutor(e *Engine, decl *parser.Decl, conn protocol.EngineConn) error {
	return conn.WriteResult(0, 0)
}
