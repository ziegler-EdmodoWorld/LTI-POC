package parser

import (
	"regexp"
)

// Token contains a token id and its lexeme
// TODO: Add Location{Path string, Line int, Offset int}
type Token struct {
	Token  int
	Lexeme string
}

// IsAWord is used to tell if a keyword token is effectively a string (i.e. a word)
func (u Token) IsAWord() bool {
	var wordRE = regexp.MustCompile("^[a-zA-Z]+$")

	return wordRE.MatchString(u.Lexeme)
}
