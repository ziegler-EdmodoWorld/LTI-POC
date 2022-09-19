package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
)

// Attribute (aka Field, Column) is a named column of a relation
type Attribute struct {
	name          string
	selectAs      string
	typeName      string
	defaultValue  interface{}
	onUpdateValue interface{}
	autoIncrement bool // TODO: rename to isAutoIncrement
	unique        bool // TODO: rename to isUnique
	isNullable    bool
}

// NewAttribute initialize a new Attribute struct
func NewAttribute(name string, typeName string, autoIncrement bool) Attribute {
	return Attribute{
		name:          name,
		typeName:      typeName,
		autoIncrement: autoIncrement,
		isNullable:    true,
	}
}

// Alias is a convenience method to get the localized, effective name (which may or may not be aliased)
func (u Attribute) Alias() string {
	if u.selectAs != "" {
		return u.selectAs
	}

	nameFields := strings.Split(u.name, ".")
	return nameFields[len(nameFields)-1]
}

// TranslateDecl traverses a Decl tree translating token sequences into Attribute settings
// TODO func (u *Attribute) TranslateDecl(decl *parser.Decl) error {...}

func parseAttribute(decl *parser.Decl) (Attribute, error) {
	attr := NewAttribute("", "", false)

	// Attribute name
	if decl.Token != parser.StringToken {
		return attr, fmt.Errorf("engine: expected attribute name, got %v", decl.Token)
	}
	attr.name = decl.Lexeme

	// Attribute type
	if len(decl.Decl) < 1 {
		return attr, fmt.Errorf("Attribute %s has no type", decl.Lexeme)
	}
	if decl.Decl[0].Token != parser.StringToken {
		return attr, fmt.Errorf("engine: expected attribute type, got %v:%v", decl.Decl[0].Token, decl.Decl[0].Lexeme)
	}
	attr.typeName = decl.Decl[0].Lexeme

	// Maybe domain and special thing like primary key
	otherDecl := decl.Decl[1:]
	for i := range otherDecl {
		log.Debug("Got %v for %s %s", otherDecl[i], attr.name, attr.typeName)
		switch otherDecl[i].Token {
		case parser.AutoincrementToken: // AUTOINCREMENT
			attr.autoIncrement = true
		case parser.UniqueToken: // UNIQUE
			attr.unique = true
		case parser.NotToken: // NOT NULL
			if len(otherDecl[i].Decl) != 1 {
				return attr, fmt.Errorf("Attribute %s has incomplete NOT NULL constraint", attr.name)
			}
			switch otherDecl[i].Decl[0].Token {
			case parser.NullToken:
				attr.isNullable = false
			}
		case parser.NullToken: // NULL
			if len(otherDecl[i].Decl) != 0 {
				return attr, fmt.Errorf("Attribute %s has NULL constraint with extra arguments", attr.name)
			}
			attr.isNullable = true
		case parser.DefaultToken: // DEFAULT <VALUE>
			log.Debug("we get a default value for %s: %s!\n", attr.name, otherDecl[i].Decl[0].Lexeme)
			switch otherDecl[i].Decl[0].Token {
			case parser.LocalTimestampToken, parser.NowToken:
				log.Debug("Setting default value to NOW() func !\n")
				attr.defaultValue = func() interface{} { return time.Now().Format(parser.DateLongFormat) }
			default:
				log.Debug("Setting default value to '%v'\n", otherDecl[i].Decl[0].Lexeme)
				attr.defaultValue = otherDecl[i].Decl[0].Lexeme
			}
		case parser.OnToken: // ON UPDATE <VALUE>
			log.Debug("we get a on update value for %s: %s!\n", attr.name, otherDecl[i].Decl[0].Decl[0].Lexeme)
			switch otherDecl[i].Decl[0].Decl[0].Token {
			case parser.LocalTimestampToken, parser.NowToken:
				log.Debug("Setting on update value to NOW() func !\n")
				attr.onUpdateValue = func() interface{} { return time.Now().Format(parser.DateLongFormat) }
			default:
				log.Debug("Setting on update value to '%v'\n", otherDecl[i].Decl[0].Decl[0].Lexeme)
				attr.onUpdateValue = otherDecl[i].Decl[0].Decl[0].Lexeme
			}
		}
	}

	if strings.ToLower(attr.typeName) == "bigserial" {
		attr.autoIncrement = true
	}

	return attr, nil
}
