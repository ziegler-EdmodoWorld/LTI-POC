package parser

import (
	"fmt"
	"strings"
)

func (p *parser) parseCreate() (*Instruction, error) {
	i := &Instruction{}

	// Set CREATE decl
	createDecl := NewDecl(p.cur())
	i.Decls = append(i.Decls, createDecl)

	// After create token, should be either
	// TABLE
	// INDEX
	// ...
	if !p.hasNext() {
		return nil, fmt.Errorf("CREATE token must be followed by TABLE, INDEX")
	}
	p.index++

	switch p.cur().Token {
	case TableToken:
		d, err := p.parseTable()
		if err != nil {
			return nil, err
		}
		createDecl.Add(d)
		break
	default:
		return nil, fmt.Errorf("Parsing error near <%s>", p.cur().Lexeme)
	}

	return i, nil
}

func (p *parser) parseTable() (*Decl, error) {
	var err error
	tableDecl := NewDecl(p.cur())
	p.index++

	// Optional: IF NOT EXISTS
	if p.is(IfToken) {
		ifDecl, err := p.consumeToken(IfToken)
		if err != nil {
			return nil, err
		}
		tableDecl.Add(ifDecl)

		if p.is(NotToken) {
			notDecl, err := p.consumeToken(NotToken)
			if err != nil {
				return nil, err
			}
			ifDecl.Add(notDecl)
			if !p.is(ExistsToken) {
				return nil, p.syntaxError()
			}
			existsDecl, err := p.consumeToken(ExistsToken)
			if err != nil {
				return nil, err
			}
			notDecl.Add(existsDecl)
		}
	}

	// Required: <TABLE-NAME>
	nameTable, err := p.parseAttribute()
	if err != nil {
		return nil, p.syntaxError()
	}
	tableDecl.Add(nameTable)

	// Required: '(' (Opening Parenthesis)
	if !p.hasNext() || p.cur().Token != BracketOpeningToken {
		return nil, fmt.Errorf("Table name token must be followed by table definition")
	}
	p.index++

	// Required: <TABLE-BODY>
	for p.index < p.tokenLen {

		// ')' (Closing parenthesis)
		if p.cur().Token == BracketClosingToken {
			p.consumeToken(BracketClosingToken)
			break
		}

		// (CONSTRAINT <CONSTRAINT-NAME>?)? ...
		if p.cur().Token == ConstraintToken {
			_, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}

			// PRIMARY KEY ( <INDEX-KEY> [, ...] )
		} else if p.cur().Token == PrimaryToken {
			_, err := p.parsePrimaryKey()
			if err != nil {
				return nil, err
			}

			// UNIQUE [INDEX | KEY] ...
		} else if p.cur().Token == UniqueToken {
			_, err := p.consumeToken(UniqueToken)
			if err != nil {
				return nil, err
			}

			_, err = p.parseTableIndex()
			if err != nil {
				return nil, err
			}

			// { INDEX | KEY } [ index_name ] [?:index_type USING { BTREE | HASH } ] '(' { col_name [ '(' length ')' ] | '(' expr ')' } [ ASC | DESC ] ',' ... ')' [?:index_option ... ]
		} else if p.cur().Token == IndexToken || p.cur().Token == KeyToken {
			_, err := p.parseTableIndex()
			if err != nil {
				return nil, err
			}

			// FOREIGN KEY ...
		} else if p.cur().Token == ForeignToken {
			_, err := p.parseTableForeignKey()
			if err != nil {
				return nil, err
			}

			// <TABLE-ATTRIBUTE>
		} else {
			// New attribute name
			newAttribute, err := p.parseQuotedToken()
			if err != nil {
				return nil, err
			}
			tableDecl.Add(newAttribute)

			newAttributeType, err := p.parseType()
			if err != nil {
				return nil, err
			}
			newAttribute.Add(newAttributeType)

			// All the following tokens until bracket or comma are column constraints.
			// Column constraints can be listed in any order.
			for p.isNot(BracketClosingToken, CommaToken) {
				switch p.cur().Token {
				case UniqueToken: // UNIQUE
					uniqueDecl, err := p.consumeToken(UniqueToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(uniqueDecl)
				case NotToken: // NOT NULL
					if _, err = p.isNext(NullToken); err == nil {
						notDecl, err := p.consumeToken(NotToken)
						if err != nil {
							return nil, err
						}
						newAttribute.Add(notDecl)
						nullDecl, err := p.consumeToken(NullToken)
						if err != nil {
							return nil, err
						}
						notDecl.Add(nullDecl)
					}
				case NullToken: // NULL
					nullDecl, err := p.consumeToken(NullToken)
					if err != nil {
						return nil, err
					}

					newAttribute.Add(nullDecl)
				case PrimaryToken: // PRIMARY KEY
					if _, err = p.isNext(KeyToken); err == nil {
						newPrimary := NewDecl(p.cur())
						newAttribute.Add(newPrimary)

						if err = p.next(); err != nil {
							return nil, fmt.Errorf("Unexpected end")
						}

						newKey := NewDecl(p.cur())
						newPrimary.Add(newKey)

						if err = p.next(); err != nil {
							return nil, fmt.Errorf("Unexpected end")
						}
					}
				case AutoincrementToken:
					autoincDecl, err := p.consumeToken(AutoincrementToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(autoincDecl)
				case WithToken: // WITH TIME ZONE
					if strings.ToLower(newAttributeType.Lexeme) == "timestamp" {
						withDecl, err := p.consumeToken(WithToken)
						if err != nil {
							return nil, err
						}
						timeDecl, err := p.consumeToken(TimeToken)
						if err != nil {
							return nil, err
						}
						zoneDecl, err := p.consumeToken(ZoneToken)
						if err != nil {
							return nil, err
						}
						newAttributeType.Add(withDecl)
						withDecl.Add(timeDecl)
						timeDecl.Add(zoneDecl)
					}
				case DefaultToken: // DEFAULT <VALUE>
					defaultDecl, err := p.consumeToken(DefaultToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(defaultDecl)
					valueDecl, err := p.consumeToken(FalseToken, StringToken, NumberToken, LocalTimestampToken, NullToken)
					if err != nil {
						return nil, err
					}
					defaultDecl.Add(valueDecl)
				case OnToken: // ON UPDATE <VALUE>
					onDecl, err := p.consumeToken(OnToken)
					if err != nil {
						return nil, err
					}

					updateDecl, err := p.consumeToken(UpdateToken)
					if err != nil {
						return nil, err
					}

					vDecl, err := p.consumeToken(FalseToken, StringToken, NumberToken, LocalTimestampToken)
					if err != nil {
						return nil, err
					}

					onDecl.Add(updateDecl)
					updateDecl.Add(vDecl)
					newAttribute.Add(onDecl)
				default:
					// Unknown column constraint
					return nil, p.syntaxError()
				}
			}
		}

		// Comma means continue to next table column
		// NOTE: With this the parser accepts ", )" and happily proceeds but this is not valid SQL (AFAIK)
		if p.cur().Token == CommaToken {
			p.index++
		}
	}

	// Optional: <TABLE-OPTIONS> - these can be listed in any order
tableOptions:
	for p.index < p.tokenLen {
		switch p.cur().Token {
		case EngineToken: // ENGINE [=] value
			engineDecl, err := p.consumeToken(EngineToken)
			if err != nil {
				return nil, err
			}

			if p.cur().Token == EqualityToken {
				if err = p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(FalseToken, StringToken, NumberToken)
			if err != nil {
				return nil, err
			}

			engineDecl.Add(vDecl)
			// TODO: tableDecl.Add(engineDecl)

		case DefaultToken: // [DEFAULT] (CHARACTER SET, CHARSET, COLLATE) [=] value
			if err := p.next(); err != nil {
				return nil, err
			}

			switch p.cur().Token {
			case CharsetToken: // CHARSET [=] value
				if err := p.next(); err != nil {
					return nil, err
				}
				charDecl := NewDecl(Token{Token: CharacterToken, Lexeme: "character"})
				setDecl := NewDecl(Token{Token: SetToken, Lexeme: "set"})

				if p.cur().Token == EqualityToken {
					if err := p.next(); err != nil {
						return nil, err
					}
				}

				vDecl, err := p.consumeToken(StringToken)
				if err != nil {
					return nil, err
				}

				charDecl.Add(setDecl)
				setDecl.Add(vDecl)
				// TODO: tableDecl.Add(charDecl)

			case CharacterToken: // CHARACTER SET [=] value
				charDecl, err := p.consumeToken(CharacterToken)
				if err != nil {
					return nil, err
				}

				setDecl, err := p.consumeToken(SetToken)
				if err != nil {
					return nil, err
				}

				if p.cur().Token == EqualityToken {
					if err := p.next(); err != nil {
						return nil, err
					}
				}

				vDecl, err := p.consumeToken(StringToken)
				if err != nil {
					return nil, err
				}

				charDecl.Add(setDecl)
				setDecl.Add(vDecl)
				// TODO: tableDecl.Add(charDecl)
			default:
				// Unknown 'table_option'
				return nil, p.syntaxError()
			}

		case CharsetToken: // CHARSET [=] value
			if err := p.next(); err != nil {
				return nil, err
			}
			charDecl := NewDecl(Token{Token: CharacterToken, Lexeme: "character"})
			setDecl := NewDecl(Token{Token: SetToken, Lexeme: "set"})

			if p.cur().Token == EqualityToken {
				if err := p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(StringToken)
			if err != nil {
				return nil, err
			}

			charDecl.Add(setDecl)
			setDecl.Add(vDecl)
			// TODO: tableDecl.Add(charDecl)

		case CharacterToken: // CHARACTER SET [=] value
			charDecl, err := p.consumeToken(CharacterToken)
			if err != nil {
				return nil, err
			}

			setDecl, err := p.consumeToken(SetToken)
			if err != nil {
				return nil, err
			}

			if p.cur().Token == EqualityToken {
				if err := p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(StringToken)
			if err != nil {
				return nil, err
			}

			charDecl.Add(setDecl)
			setDecl.Add(vDecl)
			// TODO: tableDecl.Add(charDecl)

		case SemicolonToken: // semicolon means end of instruction
			// Important NOT to consume the semicolon token

			break tableOptions

		default: // Does not appear to be a 'table_constraint' so stop processing instruction
			break tableOptions
		}
	}

	return tableDecl, nil
}

// parseTableConstraint processes tokens that should define a table constraint
// CONSTRAINT <CONSTRAINT-NAME>? ...
func (p *parser) parseTableConstraint() (*Decl, error) {
	constraintDecl, err := p.consumeToken(ConstraintToken)
	if err != nil {
		return nil, err
	}

	// Optional: <CONSTRAINT-NAME>
	if p.is(StringToken) {
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}
	}

	switch p.cur().Token {
	case PrimaryToken:
		_, err := p.parsePrimaryKey()
		if err != nil {
			return nil, err
		}
	case UniqueToken:
		_, err := p.consumeToken(UniqueToken)
		if err != nil {
			return nil, err
		}

		_, err = p.parseTableIndex()
		if err != nil {
			return nil, err
		}
	case ForeignToken:
		_, err := p.parseTableForeignKey()
		if err != nil {
			return nil, err
		}
	default:
		// Unknown constraint type
		return nil, p.syntaxError()
	}

	return constraintDecl, nil
}

func (p *parser) parsePrimaryKey() (*Decl, error) {
	primaryDecl, err := p.consumeToken(PrimaryToken)
	if err != nil {
		return nil, err
	}

	keyDecl, err := p.consumeToken(KeyToken)
	if err != nil {
		return nil, err
	}
	primaryDecl.Add(keyDecl)

	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	for {
		d, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}

		d, err = p.consumeToken(CommaToken, BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if d.Token == BracketClosingToken {
			break
		}
	}

	return primaryDecl, nil
}

// parseTableIndex processes tokens that should define a table index
// { INDEX | KEY } [ index_name ] [?:index_type USING { BTREE | HASH } ] '(' { col_name [ '(' length ')' ] | '(' expr ')' } [ ASC | DESC ] ',' ... ')' [?:index_option ... ]
func (p *parser) parseTableIndex() (*Decl, error) {
	indexDecl := NewDecl(Token{Token: IndexToken, Lexeme: "index"})

	// Required: { INDEX | KEY }
	switch p.cur().Token {
	case IndexToken, KeyToken:
		_, err := p.consumeToken(IndexToken, KeyToken)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Table INDEX definition must start with INDEX or KEY")
	}

	// Optional: <INDEX-NAME>
	if p.is(StringToken) {
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}
	}

	// Optional: <INDEX-TYPE> := USING { BTREE | HASH }
	if p.is(UsingToken) {
		_, err := p.consumeToken(UsingToken)
		if err != nil {
			return nil, err
		}

		if p.is(BtreeToken) {
			_, err := p.consumeToken(BtreeToken)
			if err != nil {
				return nil, err
			}
		} else if p.is(HashToken) {
			_, err := p.consumeToken(HashToken)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, p.syntaxError()
		}
	}

	// Required: '('
	_, err := p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <INDEX-KEY> [ ASC | DESC ] [, <INDEX-KEY> [ ASC | DESC ] ]* ')'
	for {
		// Required: <INDEX-KEY>
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}

		// Optional: 'ASC' | 'DESC'
		if p.is(AscToken) {
			_, err := p.consumeToken(AscToken)
			if err != nil {
				return nil, err
			}
		} else if p.is(DescToken) {
			_, err := p.consumeToken(DescToken)
			if err != nil {
				return nil, err
			}
		}

		// Optional: ',' | ')'
		n, err := p.consumeToken(CommaToken, BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == BracketClosingToken {
			break
		}
	}

	return indexDecl, nil
}

// parseTableForeignKey processes tokens that should define a table foreign key
// FOREIGN KEY ...
func (p *parser) parseTableForeignKey() (*Decl, error) {
	// Required: FOREIGN
	foreignDecl, err := p.consumeToken(ForeignToken)
	if err != nil {
		return nil, err
	}

	// Required: KEY
	keyDecl, err := p.consumeToken(KeyToken)
	if err != nil {
		return nil, err
	}
	foreignDecl.Add(keyDecl)

	// Optional: <FK-NAME>
	if p.is(StringToken) {
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}
	}

	// Required: '('
	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <FK-INDEX> [, <FK-INDEX>]* ')'
	for {
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(CommaToken, BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == BracketClosingToken {
			break
		}
	}

	// Optional: REFERENCES ...
	if p.is(ReferencesToken) {
		_, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
	}

	return foreignDecl, nil
}

// parseTableReference processes tokens that should define a table reference
// REFERENCES ...
func (p *parser) parseTableReference() (*Decl, error) {
	// Required: REFERENCES
	referencesDecl, err := p.consumeToken(ReferencesToken)
	if err != nil {
		return nil, err
	}

	// Required: <TABLE-NAME>
	_, err = p.consumeToken(StringToken)
	if err != nil {
		return nil, err
	}

	// Required: '('
	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <KEY-PART> [, <KEY-PART>]* ')'
	for {
		_, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(CommaToken, BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == BracketClosingToken {
			break
		}
	}

	// Optional: MATCH ...
	if p.is(MatchToken) {
		_, err := p.consumeToken(MatchToken)
		if err != nil {
			return nil, err
		}

		switch p.cur().Token {
		case FullToken:
			_, err := p.consumeToken(FullToken)
			if err != nil {
				return nil, err
			}
		case PartialToken:
			_, err := p.consumeToken(PartialToken)
			if err != nil {
				return nil, err
			}
		case SimpleToken:
			_, err := p.consumeToken(SimpleToken)
			if err != nil {
				return nil, err
			}
		default:
			// Unknown match type
			return nil, p.syntaxError()
		}
	}

	// Optional: {ON ...}+
	if p.is(OnToken) {
		for {
			// Required: ON
			_, err := p.consumeToken(OnToken)
			if err != nil {
				return nil, err
			}

			// Required: (UPDATE | DELETE) <REFERENCE-OPTION>
			switch p.cur().Token {
			case UpdateToken:
				_, err := p.consumeToken(UpdateToken)
				if err != nil {
					return nil, err
				}

				_, err = p.parseTableReferenceOption()
				if err != nil {
					return nil, err
				}
			case DeleteToken:
				_, err := p.consumeToken(DeleteToken)
				if err != nil {
					return nil, err
				}

				_, err = p.parseTableReferenceOption()
				if err != nil {
					return nil, err
				}
			default:
				// Unknown on reference option type
				return nil, p.syntaxError()
			}

			// Repeat
			if !p.is(OnToken) {
				break
			}
		}
	}

	return referencesDecl, nil
}

// parseTableReferenceOption processes tokens that should define a table reference option
func (p *parser) parseTableReferenceOption() (*Decl, error) {
	switch p.cur().Token {
	case RestrictToken:
		d, err := p.consumeToken(RestrictToken)
		if err != nil {
			return nil, err
		}
		return d, nil
	case CascadeToken:
		d, err := p.consumeToken(CascadeToken)
		if err != nil {
			return nil, err
		}
		return d, nil
	case SetToken:
		d, err := p.consumeToken(SetToken)
		if err != nil {
			return nil, err
		}

		switch p.cur().Token {
		case NullToken:
			n, err := p.consumeToken(NullToken)
			if err != nil {
				return nil, err
			}
			d.Add(n)
			return d, nil
		case DefaultToken:
			n, err := p.consumeToken(DefaultToken)
			if err != nil {
				return nil, err
			}
			d.Add(n)
			return d, nil
		default:
			// Unknown option
			return nil, p.syntaxError()
		}
	case NoToken:
		d, err := p.consumeToken(NoToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(ActionToken)
		if err != nil {
			return nil, err
		}
		d.Add(n)
		return d, nil
	default:
		// Unknown option type
		return nil, p.syntaxError()
	}
}
