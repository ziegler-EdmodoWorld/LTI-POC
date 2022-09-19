package engine

import (
	"fmt"

	"github.com/mlhoyt/ramsql/engine/protocol"
)

func truncateTable(e *Engine, table *Table, conn protocol.EngineConn) error {
	var rowsDeleted int64

	// get relations and write lock them
	r := e.relation(table.name)
	if r == nil {
		return fmt.Errorf("Table %v not found", table.name)
	}
	r.Lock()
	defer r.Unlock()

	if r.rows != nil {
		rowsDeleted = int64(len(r.rows))
	}
	r.rows = make([]*Tuple, 0)

	return conn.WriteResult(0, rowsDeleted)
}
