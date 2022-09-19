package engine

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

type executor func(*Engine, *parser.Decl, protocol.EngineConn) error

// Engine is the root struct of RamSQL server
type Engine struct {
	endpoint     protocol.EngineEndpoint
	relations    map[string]*Relation
	opsExecutors map[int]executor

	// Any value send to this channel (through Engine.stop)
	// Will stop the listening loop
	stop chan bool

	sync.Mutex
}

// New initialize a new RamSQL server
func New(endpoint protocol.EngineEndpoint) (e *Engine, err error) {

	e = &Engine{
		endpoint: endpoint,
	}

	e.stop = make(chan bool)

	e.opsExecutors = map[int]executor{
		parser.CreateToken:   createExecutor,
		parser.DeleteToken:   deleteExecutor,
		parser.DropToken:     dropExecutor,
		parser.ExistsToken:   existsExecutor,
		parser.GrantToken:    grantExecutor,
		parser.IfToken:       ifExecutor,
		parser.InsertToken:   insertIntoTableExecutor,
		parser.NotToken:      notExecutor,
		parser.SelectToken:   selectExecutor,
		parser.TableToken:    createTableExecutor,
		parser.TruncateToken: truncateExecutor,
		parser.UpdateToken:   updateExecutor,
	}

	e.relations = make(map[string]*Relation)

	err = e.start()
	if err != nil {
		return nil, err
	}

	return
}

func (e *Engine) start() (err error) {
	go e.listen()
	return nil
}

// Stop shutdown the RamSQL server
func (e *Engine) Stop() {
	if e.stop == nil {
		// already stopped
		return
	}
	go func() {
		e.stop <- true
		close(e.stop)
		e.stop = nil
	}()
}

func (e *Engine) relation(name string) *Relation {
	// Lock ?
	r := e.relations[name]
	// Unlock ?

	return r
}

func (e *Engine) drop(name string) {
	e.Lock()
	delete(e.relations, name)
	e.Unlock()
}

func (e *Engine) listen() {
	newConnectionChannel := make(chan protocol.EngineConn)

	go func() {
		for {
			conn, err := e.endpoint.Accept()
			if err != nil {
				e.Stop()
				return
			}

			newConnectionChannel <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnectionChannel:
			go e.handleConnection(conn)
			break

		case <-e.stop:
			e.endpoint.Close()
			return
		}
	}

}

func (e *Engine) handleConnection(conn protocol.EngineConn) {

	for {
		stmt, err := conn.ReadStatement()
		if err == io.EOF {
			// Todo: close engine if there is no conn left
			return
		}
		if err != nil {
			log.Warning("Enginge.handleConnection: cannot read : %s", err)
			return
		}

		instructions, err := parser.ParseInstruction(stmt)
		if err != nil {
			conn.WriteError(err)
			continue
		}

		err = e.executeQueries(instructions, conn)
		if err != nil {
			conn.WriteError(err)
			continue
		}
	}
}

func (e *Engine) executeQueries(instructions []parser.Instruction, conn protocol.EngineConn) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("fatal error: %s", r)
			return
		}
	}()

	for _, i := range instructions {
		err = e.executeQuery(i, conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) executeQuery(i parser.Instruction, conn protocol.EngineConn) error {

	if e.opsExecutors[i.Decls[0].Token] != nil {
		return e.opsExecutors[i.Decls[0].Token](e, i.Decls[0], conn)
	}

	return errors.New("Not Implemented")
}
