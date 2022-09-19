package engine

import (
	"fmt"
	"strconv"

	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

/*
|-> SELECT
	|-> *
	|-> FROM
		|-> account
	|-> WHERE
		|-> email
			|-> =
			|-> foo@bar.com
*/
func selectExecutor(e *Engine, selectDecl *parser.Decl, conn protocol.EngineConn) error {
	var attributes []Attribute
	var tables []*Table
	var predicates []PredicateLinker
	var functors []selectFunctor
	var joiners []joiner
	var err error

	selectDecl.Stringy(0)
	for i := range selectDecl.Decl {
		switch selectDecl.Decl[i].Token {
		case parser.FromToken:
			// get selected tables
			tables = fromExecutor(selectDecl.Decl[i])
		case parser.WhereToken:
			// get WHERE declaration
			pred, err := whereExecutor2(e, selectDecl.Decl[i].Decl, tables[0].name)
			if err != nil {
				return err
			}
			predicates = []PredicateLinker{pred}
		case parser.JoinToken:
			j, err := joinExecutor(selectDecl.Decl[i])
			if err != nil {
				return err
			}
			joiners = append(joiners, j)
		case parser.OrderToken:
			orderFunctor, err := orderbyExecutor(selectDecl.Decl[i], tables)
			if err != nil {
				return err
			}
			functors = append(functors, orderFunctor)
		case parser.LimitToken:
			limit, err := strconv.Atoi(selectDecl.Decl[i].Decl[0].Lexeme)
			if err != nil {
				return fmt.Errorf("wrong limit value: %s", err)
			}
			conn = limitedConn(conn, limit)
		case parser.OffsetToken:
			offset, err := strconv.Atoi(selectDecl.Decl[i].Decl[0].Lexeme)
			if err != nil {
				return fmt.Errorf("wrong offset value: %s", err)
			}
			conn = offsetedConn(conn, offset)
		}
	}

	for i := range selectDecl.Decl {
		if selectDecl.Decl[i].Token != parser.StringToken &&
			selectDecl.Decl[i].Token != parser.StarToken &&
			selectDecl.Decl[i].Token != parser.CountToken {
			continue
		}

		// get attribute to selected
		attr, err := getSelectedAttribute(e, selectDecl.Decl[i], tables)
		if err != nil {
			return err
		}
		attributes = append(attributes, attr...)

	}

	if len(functors) == 0 {
		// Instanciate a new select functor
		functors, err = getSelectFunctors(selectDecl)
		if err != nil {
			return err
		}
	}

	err = generateVirtualRows(e, attributes, conn, tables[0].name, joiners, predicates, functors)
	if err != nil {
		return err
	}

	return nil
}
