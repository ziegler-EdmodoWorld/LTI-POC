package parser

import (
	"github.com/mlhoyt/ramsql/engine/log"
)

func (p *parser) parseDrop() (*Instruction, error) {
	i := &Instruction{}

	// Required: DROP
	dropDecl, err := p.consumeToken(DropToken)
	if err != nil {
		log.Debug("WTF\n")
		return nil, err
	}
	i.Decls = append(i.Decls, dropDecl)

	// Required: TABLE
	tableDecl, err := p.consumeToken(TableToken)
	if err != nil {
		log.Debug("Consume table !\n")
		return nil, err
	}
	dropDecl.Add(tableDecl)

	// Optional: IF EXISTS
	if p.is(IfToken) {
		ifDecl, err := p.consumeToken(IfToken)
		if err != nil {
			return nil, err
		}
		tableDecl.Add(ifDecl)

		// Required: EXISTS
		if !p.is(ExistsToken) {
			return nil, p.syntaxError()
		}

		existsDecl, err := p.consumeToken(ExistsToken)
		if err != nil {
			return nil, err
		}
		ifDecl.Add(existsDecl)
	}

	// Required: <TABLE-NAME>
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		log.Debug("UH ?\n")
		return nil, err
	}
	tableDecl.Add(nameDecl)

	return i, nil
}
