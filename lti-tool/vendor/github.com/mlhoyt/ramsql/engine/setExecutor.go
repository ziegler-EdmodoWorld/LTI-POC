package engine

import (
	"github.com/mlhoyt/ramsql/engine/parser"
)

/*
	|-> set
	      |-> email
					|-> =
					|-> roger@gmail.com
*/
func setExecutor(setDecl *parser.Decl) (map[string]interface{}, error) {

	values := make(map[string]interface{})

	for _, attr := range setDecl.Decl {
		values[attr.Lexeme] = attr.Decl[1].Lexeme
	}

	return values, nil
}
