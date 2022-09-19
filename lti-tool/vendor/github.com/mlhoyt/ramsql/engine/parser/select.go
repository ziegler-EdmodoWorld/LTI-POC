package parser

import (
	"fmt"
)

func (p *parser) parseSelect() (*Instruction, error) {
	i := &Instruction{}
	var err error

	// Create select decl
	selectDecl := NewDecl(p.cur())
	i.Decls = append(i.Decls, selectDecl)

	// After select token, should be either
	// a StarToken
	// a list of table names + (StarToken Or Attribute)
	// a builtin func (COUNT, MAX, ...)
	if err = p.next(); err != nil {
		return nil, fmt.Errorf("SELECT token must be followed by attributes to select")
	}

	for {
		if p.is(CountToken) {
			attrDecl, err := p.parseBuiltinFunc()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		} else {
			attrDecl, err := p.parseAttribute()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		}

		// If comma, loop again.
		if p.is(CommaToken) {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	// Required: FROM
	if p.cur().Token != FromToken {
		return nil, fmt.Errorf("Syntax error near %v", p.cur())
	}
	fromDecl := NewDecl(p.cur())
	selectDecl.Add(fromDecl)

	// Now must be a list of table
	for {
		// string
		if err = p.next(); err != nil {
			return nil, fmt.Errorf("Unexpected end. Syntax error near %v", p.cur())
		}
		tableNameDecl, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		fromDecl.Add(tableNameDecl)

		// If no next, then it's implicit where
		if !p.hasNext() {
			addImplicitWhereAll(selectDecl)
			return i, nil
		}
		// if not comma, break
		if p.cur().Token != CommaToken {
			break // No more table
		}
	}

	// (INNER | ((LEFT|RIGHT) (OUTER)?))? JOIN
	for p.is(InnerToken, LeftToken, RightToken, OuterToken, JoinToken) {
		// Optional: INNER
		// var innerJoinDecl *Decl
		if p.is(InnerToken) {
			_, err := p.consumeToken(InnerToken)
			if err != nil {
				return nil, err
			}
		}

		// Optional: LEFT, RIGHT
		// var dirOuterJoinDecl *Decl
		if p.is(LeftToken) {
			_, err := p.consumeToken(LeftToken)
			if err != nil {
				return nil, err
			}
		} else if p.is(RightToken) {
			_, err := p.consumeToken(RightToken)
			if err != nil {
				return nil, err
			}
		}

		// Optional: OUTER
		// var outerJoinDecl *Decl
		if p.is(OuterToken) {
			_, err := p.consumeToken(OuterToken)
			if err != nil {
				return nil, err
			}
		}

		if !p.is(JoinToken) {
			return nil, fmt.Errorf("Syntax error near %v.  Expected JOIN", p.cur())
		}
		joinDecl, err := p.parseJoin()
		if err != nil {
			return nil, err
		}
		// FIXME: Need to annotate joinDecl with innerJoinDecl or outerJoinDecl (and dirOuterJoinDecl)
		selectDecl.Add(joinDecl)
	}

	// Optional: WHERE ..., ORDER [BY] ..., LIMIT ..., OFFSET ..., FOR ...
	hazWhereClause := false
	for {
		switch p.cur().Token {
		case WhereToken:
			err := p.parseWhere(selectDecl)
			if err != nil {
				return nil, err
			}
			hazWhereClause = true
		case OrderToken:
			if hazWhereClause == false {
				// WHERE clause is implicit
				addImplicitWhereAll(selectDecl)
			}
			err := p.parseOrderBy(selectDecl)
			if err != nil {
				return nil, err
			}
		case LimitToken:
			limitDecl, err := p.consumeToken(LimitToken)
			if err != nil {
				return nil, err
			}
			selectDecl.Add(limitDecl)
			numDecl, err := p.consumeToken(NumberToken)
			if err != nil {
				return nil, err
			}
			limitDecl.Add(numDecl)
		case OffsetToken:
			offsetDecl, err := p.consumeToken(OffsetToken)
			if err != nil {
				return nil, err
			}
			selectDecl.Add(offsetDecl)
			offsetValue, err := p.consumeToken(NumberToken)
			if err != nil {
				return nil, err
			}
			offsetDecl.Add(offsetValue)
		case ForToken:
			err := p.parseForUpdate(selectDecl)
			if err != nil {
				return nil, err
			}
		default:
			return i, nil
		}
	}
}

func addImplicitWhereAll(decl *Decl) {

	whereDecl := &Decl{
		Token:  WhereToken,
		Lexeme: "where",
	}
	whereDecl.Add(&Decl{
		Token:  NumberToken,
		Lexeme: "1",
	})

	decl.Add(whereDecl)
}

func (p *parser) parseForUpdate(decl *Decl) error {
	// Optionnal
	if !p.is(ForToken) {
		return nil
	}

	d, err := p.consumeToken(ForToken)
	if err != nil {
		return err
	}

	u, err := p.consumeToken(UpdateToken)
	if err != nil {
		return err
	}

	d.Add(u)
	decl.Add(d)
	return nil
}
