package parser

import (
	"fmt"
	"unicode"

	"github.com/mlhoyt/ramsql/engine/log"
)

type lexer struct {
	tokens         []Token
	instruction    []byte
	instructionLen int
	pos            int
}

// SQL Tokens
const (
	ActionToken         = iota // Second-order
	AndToken                   // Second-order
	AsToken                    // Second-order
	AscToken                   // Second-order
	AutoincrementToken         // Second-order
	BacktickToken              // Punctuation
	BracketClosingToken        // Punctuation
	BracketOpeningToken        // Punctuation
	BtreeToken                 // Second-order
	ByToken                    // Second-order
	CascadeToken               // Second-order
	CharacterToken             // Second-order
	CharsetToken               // Second-order
	CommaToken                 // Punctuation
	ConstraintToken            // Second-order
	CountToken                 // Second-order
	CreateToken                // First-order
	DateToken                  // Type
	DefaultToken               // Second-order
	DeleteToken                // First-order
	DescToken                  // Second-order
	DoubleQuoteToken           // Quote
	DropToken                  // First-order
	EngineToken                // Second-order
	EqualityToken              // Quote
	ExistsToken                // Second-order
	ExplainToken               // First-order
	FalseToken                 // Second-order
	ForToken                   // Second-order
	ForeignToken               // Second-order
	FromToken                  // Second-order
	FullToken                  // Second-order
	GrantToken                 // First-order
	GreaterOrEqualToken        // Punctuation
	HashToken                  // Second-order
	IfToken                    // Second-order
	InToken                    // Second-order
	IndexToken                 // Second-order
	InnerToken                 // Second-order
	InsertToken                // First-order
	IntToken                   // Type
	IntoToken                  // Second-order
	IsToken                    // Second-order
	JoinToken                  // Second-order
	KeyToken                   // Type
	LeftToken                  // Second-order
	LeftDipleToken             // Punctuation
	LessOrEqualToken           // Punctuation
	LimitToken                 // Second-order
	LocalTimestampToken        // Second-order
	MatchToken                 // Second-order
	NoToken                    // Second-order
	NotToken                   // Second-order
	NowToken                   // Second-order
	NullToken                  // Second-order
	NumberToken                // Type
	OffsetToken                // Second-order
	OnToken                    // Second-order
	OrToken                    // Second-order
	OrderToken                 // Second-order
	OuterToken                 // Second-order
	PartialToken               // Quote
	PeriodToken                // Quote
	PrimaryToken               // Type
	ReferencesToken            // Second-order
	ReturningToken             // Second-order
	RestrictToken              // Second-order
	RightToken                 // Second-order
	RightDipleToken            // Punctuation
	SelectToken                // First-order
	SemicolonToken             // Punctuation
	SetToken                   // Second-order
	SimpleToken                // Second-order
	SimpleQuoteToken           // Quote
	SpaceToken                 // Punctuation
	StarToken                  // Quote
	StringToken                // Type
	TableToken                 // Second-order
	TextToken                  // Type
	TimeToken                  // Second-order
	TrueToken                  // Second-order
	TruncateToken              // First-order
	UniqueToken                // Second-order
	UpdateToken                // First-order
	UsingToken                 // Second-order
	ValuesToken                // Second-order
	WhereToken                 // Second-order
	WithToken                  // Second-order
	ZoneToken                  // Second-order
)

// Matcher tries to match given string to an SQL token
type Matcher func() bool

//go:generate go run ../../utils/lexer-generate-matcher.go --init
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "action"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "and"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "as"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "asc"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "autoincrement" --lexeme "auto_increment"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "`" --name Backtick
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme ")" --name BracketClosing
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "(" --name BracketOpening
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "btree"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "by"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "cascade"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "character"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "charset"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "," --name Comma
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "constraint"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "count"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "create"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "default"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "delete"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "desc"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "drop"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "engine"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "=" --name Equality
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "exists"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "false"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "for"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "foreign"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "from"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "full"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "grant"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme ">=" --name GreaterOrEqual
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "hash"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "if"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "in"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "index"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "inner"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "insert"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "into"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "is"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "join"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "key"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "left"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "<" --name LeftDiple
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "<=" --name LessOrEqual
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "limit"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "localtimestamp" --lexeme "current_timestamp" --name LocalTimestamp
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "match"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "no"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "not"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "now()" --name Now
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "null"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "offset"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "on"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "or"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "order"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "outer"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "partial"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "." --name Period
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "primary"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "references"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "restrict"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "returning"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "right"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme ">" --name RightDiple
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "select"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme ";" --name Semicolon
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "set"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "simple"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "*" --name Star
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "table"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "time"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "true"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "truncate"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "unique"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "update"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "using"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "values"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "where"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "with"
//go:generate go run ../../utils/lexer-generate-matcher.go --lexeme "zone"

func (l *lexer) lex(instruction []byte) ([]Token, error) {
	l.instructionLen = len(instruction)
	l.tokens = nil
	l.instruction = instruction
	l.pos = 0
	securityPos := 0

	var matchers []Matcher
	// Punctuation Matcher
	matchers = append(matchers, l.MatchSpaceToken)
	matchers = append(matchers, l.MatchSemicolonToken)
	matchers = append(matchers, l.MatchCommaToken)
	matchers = append(matchers, l.MatchBracketOpeningToken)
	matchers = append(matchers, l.MatchBracketClosingToken)
	matchers = append(matchers, l.MatchStarToken)
	matchers = append(matchers, l.MatchSimpleQuoteToken)
	matchers = append(matchers, l.MatchEqualityToken)
	matchers = append(matchers, l.MatchPeriodToken)
	matchers = append(matchers, l.MatchDoubleQuoteToken)
	matchers = append(matchers, l.MatchLessOrEqualToken)
	matchers = append(matchers, l.MatchLeftDipleToken)
	matchers = append(matchers, l.MatchGreaterOrEqualToken)
	matchers = append(matchers, l.MatchRightDipleToken)
	matchers = append(matchers, l.MatchBacktickToken)
	// First order Matcher
	matchers = append(matchers, l.MatchCreateToken)
	matchers = append(matchers, l.MatchDeleteToken)
	matchers = append(matchers, l.MatchDropToken)
	matchers = append(matchers, l.MatchGrantToken)
	matchers = append(matchers, l.MatchInsertToken)
	matchers = append(matchers, l.MatchSelectToken)
	matchers = append(matchers, l.MatchTruncateToken)
	matchers = append(matchers, l.MatchUpdateToken)
	// Second order Matcher
	matchers = append(matchers, l.MatchActionToken)
	matchers = append(matchers, l.MatchAndToken)
	matchers = append(matchers, l.MatchAscToken)
	matchers = append(matchers, l.MatchAsToken)
	matchers = append(matchers, l.MatchAutoincrementToken)
	matchers = append(matchers, l.MatchBtreeToken)
	matchers = append(matchers, l.MatchByToken)
	matchers = append(matchers, l.MatchCascadeToken)
	matchers = append(matchers, l.MatchCharacterToken)
	matchers = append(matchers, l.MatchCharsetToken)
	matchers = append(matchers, l.MatchConstraintToken)
	matchers = append(matchers, l.MatchCountToken)
	matchers = append(matchers, l.MatchDefaultToken)
	matchers = append(matchers, l.MatchDescToken)
	matchers = append(matchers, l.MatchEngineToken)
	matchers = append(matchers, l.MatchExistsToken)
	matchers = append(matchers, l.MatchFalseToken)
	matchers = append(matchers, l.MatchForeignToken)
	matchers = append(matchers, l.MatchForToken)
	matchers = append(matchers, l.MatchFromToken)
	matchers = append(matchers, l.MatchFullToken)
	matchers = append(matchers, l.MatchHashToken)
	matchers = append(matchers, l.MatchIfToken)
	matchers = append(matchers, l.MatchIndexToken)
	matchers = append(matchers, l.MatchInnerToken)
	matchers = append(matchers, l.MatchIntoToken)
	matchers = append(matchers, l.MatchInToken)
	matchers = append(matchers, l.MatchIsToken)
	matchers = append(matchers, l.MatchJoinToken)
	matchers = append(matchers, l.MatchKeyToken)
	matchers = append(matchers, l.MatchLeftToken)
	matchers = append(matchers, l.MatchLimitToken)
	matchers = append(matchers, l.MatchLocalTimestampToken)
	matchers = append(matchers, l.MatchMatchToken)
	matchers = append(matchers, l.MatchNotToken)
	matchers = append(matchers, l.MatchNowToken)
	matchers = append(matchers, l.MatchNoToken)
	matchers = append(matchers, l.MatchNullToken)
	matchers = append(matchers, l.MatchOffsetToken)
	matchers = append(matchers, l.MatchOnToken)
	matchers = append(matchers, l.MatchOrderToken)
	matchers = append(matchers, l.MatchOrToken)
	matchers = append(matchers, l.MatchOuterToken)
	matchers = append(matchers, l.MatchPartialToken)
	matchers = append(matchers, l.MatchPrimaryToken)
	matchers = append(matchers, l.MatchReferencesToken)
	matchers = append(matchers, l.MatchRestrictToken)
	matchers = append(matchers, l.MatchReturningToken)
	matchers = append(matchers, l.MatchRightToken)
	matchers = append(matchers, l.MatchSetToken)
	matchers = append(matchers, l.MatchSimpleToken)
	matchers = append(matchers, l.MatchTableToken)
	matchers = append(matchers, l.MatchTimeToken)
	matchers = append(matchers, l.MatchUniqueToken)
	matchers = append(matchers, l.MatchUsingToken)
	matchers = append(matchers, l.MatchValuesToken)
	matchers = append(matchers, l.MatchWhereToken)
	matchers = append(matchers, l.MatchWithToken)
	matchers = append(matchers, l.MatchZoneToken)
	// Type Matcher
	matchers = append(matchers, l.MatchEscapedStringToken)
	matchers = append(matchers, l.MatchDateToken)
	matchers = append(matchers, l.MatchNumberToken)
	matchers = append(matchers, l.MatchStringToken)

	var r bool
	for l.pos < l.instructionLen {
		// fmt.Printf("Tokens : %v\n\n", l.tokens)

		r = false
		for _, m := range matchers {
			if r = m(); r == true {
				securityPos = l.pos
				break
			}
		}

		if r {
			continue
		}

		if l.pos == securityPos {
			log.Warning("Cannot lex <%s>, stuck at pos %d -> [%c]", l.instruction, l.pos, l.instruction[l.pos])
			return nil, fmt.Errorf("Cannot lex instruction. Syntax error near %s", instruction[l.pos:])
		}
		securityPos = l.pos
	}

	return l.tokens, nil
}

func (l *lexer) MatchSpaceToken() bool {

	if unicode.IsSpace(rune(l.instruction[l.pos])) {
		t := Token{
			Token:  SpaceToken,
			Lexeme: " ",
		}
		l.tokens = append(l.tokens, t)
		l.pos++
		return true
	}

	return false
}

func (l *lexer) MatchStringToken() bool {

	i := l.pos
	for i < l.instructionLen &&
		(unicode.IsLetter(rune(l.instruction[i])) ||
			unicode.IsDigit(rune(l.instruction[i])) ||
			l.instruction[i] == '_' ||
			l.instruction[i] == '@' /* || l.instruction[i] == '.'*/) {
		i++
	}

	if i != l.pos {
		t := Token{
			Token:  StringToken,
			Lexeme: string(l.instruction[l.pos:i]),
		}
		l.tokens = append(l.tokens, t)
		l.pos = i
		return true
	}

	return false
}

func (l *lexer) MatchNumberToken() bool {

	i := l.pos
	for i < l.instructionLen && unicode.IsDigit(rune(l.instruction[i])) {
		i++
	}

	if i != l.pos {
		t := Token{
			Token:  NumberToken,
			Lexeme: string(l.instruction[l.pos:i]),
		}
		l.tokens = append(l.tokens, t)
		l.pos = i
		return true
	}

	return false
}

// MatchDateToken prefers time.RFC3339Nano but will match a few others as well
func (l *lexer) MatchDateToken() bool {

	i := l.pos
	for i < l.instructionLen &&
		l.instruction[i] != ',' &&
		l.instruction[i] != ')' {
		i++
	}

	data := string(l.instruction[l.pos:i])

	_, err := ParseDate(data)
	if err != nil {
		return false
	}

	t := Token{
		Token:  StringToken,
		Lexeme: data,
	}

	l.tokens = append(l.tokens, t)
	l.pos = i
	return true
}

func (l *lexer) MatchDoubleQuoteToken() bool {

	if l.instruction[l.pos] == '"' {

		t := Token{
			Token:  DoubleQuoteToken,
			Lexeme: "\"",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.MatchDoubleQuotedStringToken() {
			t := Token{
				Token:  DoubleQuoteToken,
				Lexeme: "\"",
			}
			l.tokens = append(l.tokens, t)
			l.pos++
			return true
		}

		return true
	}

	return false
}

func (l *lexer) MatchEscapedStringToken() bool {
	i := l.pos
	if l.instruction[i] != '$' || l.instruction[i+1] != '$' {
		return false
	}
	i += 2

	for i+1 < l.instructionLen && !(l.instruction[i] == '$' && l.instruction[i+1] == '$') {
		i++
	}
	i++

	if i == l.instructionLen {
		return false
	}

	tok := NumberToken
	escaped := l.instruction[l.pos+2 : i-1]

	for _, r := range escaped {
		if unicode.IsDigit(rune(r)) == false {
			tok = StringToken
		}
	}

	_, err := ParseDate(string(escaped))
	if err == nil {
		tok = DateToken
	}

	t := Token{
		Token:  tok,
		Lexeme: string(escaped),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i + 1

	return true
}

func (l *lexer) MatchDoubleQuotedStringToken() bool {
	i := l.pos
	for i < l.instructionLen && l.instruction[i] != '"' {
		i++
	}

	t := Token{
		Token:  StringToken,
		Lexeme: string(l.instruction[l.pos:i]),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i

	return true
}

func (l *lexer) MatchSimpleQuoteToken() bool {

	if l.instruction[l.pos] == '\'' {

		t := Token{
			Token:  SimpleQuoteToken,
			Lexeme: "'",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.MatchSingleQuotedStringToken() {
			t := Token{
				Token:  SimpleQuoteToken,
				Lexeme: "'",
			}
			l.tokens = append(l.tokens, t)
			l.pos++
			return true
		}

		return true
	}

	return false
}

func (l *lexer) MatchSingleQuotedStringToken() bool {
	i := l.pos
	for i < l.instructionLen && l.instruction[i] != '\'' {
		i++
	}

	t := Token{
		Token:  StringToken,
		Lexeme: string(l.instruction[l.pos:i]),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i

	return true
}

func (l *lexer) MatchSingle(char byte, token int) bool {

	if l.pos > l.instructionLen {
		return false
	}

	if l.instruction[l.pos] != char {
		return false
	}

	t := Token{
		Token:  token,
		Lexeme: string(char),
	}

	l.tokens = append(l.tokens, t)
	l.pos++
	return true
}

func (l *lexer) Match(str []byte, token int) bool {

	if l.pos+len(str)-1 > l.instructionLen {
		return false
	}

	// Check for lowercase and uppercase
	for i := range str {
		if unicode.ToLower(rune(l.instruction[l.pos+i])) != unicode.ToLower(rune(str[i])) {
			return false
		}
	}

	// if next character is still a string, it means it doesn't match
	// ie: COUNT shoulnd match COUNTRY
	if l.instructionLen > l.pos+len(str) {
		if unicode.IsLetter(rune(l.instruction[l.pos+len(str)])) ||
			l.instruction[l.pos+len(str)] == '_' {
			return false
		}
	}

	t := Token{
		Token:  token,
		Lexeme: string(str),
	}

	l.tokens = append(l.tokens, t)
	l.pos += len(t.Lexeme)
	return true
}
