package engine

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
)

/*
|-> INTO
    |-> user
        |-> last_name
        |-> first_name
        |-> email
*/
func getRelation(e *Engine, intoDecl *parser.Decl) (*Relation, []*parser.Decl, error) {

	// Decl[0] is the table name
	r := e.relation(intoDecl.Decl[0].Lexeme)
	if r == nil {
		return nil, nil, errors.New("table " + intoDecl.Decl[0].Lexeme + " does not exists")
	}

	for i := range intoDecl.Decl[0].Decl {
		err := attributeExistsInTable(e, intoDecl.Decl[0].Decl[i].Lexeme, intoDecl.Decl[0].Lexeme)
		if err != nil {
			return nil, nil, err
		}
	}

	return r, intoDecl.Decl[0].Decl, nil
}

func insert(r *Relation, attributes []*parser.Decl, values []*parser.Decl, returnedID string) (int64, error) {
	var assigned = false
	var id int64
	var valuesindex int

	// Create tuple
	t := NewTuple()
	for attrindex, attr := range r.table.attributes {
		assigned = false

		for x, decl := range attributes {

			if attr.name == decl.Lexeme && attr.autoIncrement == false {
				// Before adding value in tuple, check it's not a builtin func or arithmetic operation
				switch values[x].Token {
				case parser.NowToken:
					t.Append(time.Now().Format(parser.DateLongFormat))
				default:
					t.Append(values[x].Lexeme)

				}
				valuesindex = x
				assigned = true

				if returnedID == attr.name {
					var err error
					id, err = strconv.ParseInt(values[x].Lexeme, 10, 64)
					if err != nil {
						return 0, err
					}
				}
			}
		}

		// If attribute is AUTO INCREMENT then compute and assign it
		if attr.autoIncrement {
			id = int64(len(r.rows) + 1)
			t.Append(id)

			assigned = true
		}

		// If attribute is UNIQUE then validate it is so
		if attr.unique {
			for i := range r.rows { // check all value already in relation (yup, no index tree)
				if r.rows[i].Values[attrindex].(string) == string(values[valuesindex].Lexeme) {
					return 0, fmt.Errorf("UNIQUE constraint violation")
				}
			}
		}

		// If value was not explictly set (or implicitly computed) then use the default value
		if assigned == false {
			switch val := attr.defaultValue.(type) {
			case func() interface{}:
				v := (func() interface{})(val)()
				log.Debug("Setting func value '%v' to %s\n", v, attr.name)
				t.Append(v)
			default:
				if val == nil && !attr.isNullable {
					return 0, fmt.Errorf("Field '%s' with constraint 'NOT NULL' doesn't have a default value", attr.name)
				}
				log.Debug("Setting default value '%v' to %s\n", val, attr.name)
				t.Append(attr.defaultValue)
			}
		}
	}

	log.Info("New tuple : %v", t)

	// Insert tuple
	err := r.Insert(t)
	if err != nil {
		return 0, err
	}

	return id, nil
}
