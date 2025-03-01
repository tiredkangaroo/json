package token

import (
	"fmt"
)

type Type uint8

const (
	LBRACKET Type = iota // {
	RBRACKET             // }

	LBRACE // [
	RBRACE // ]

	LITERAL // "hello"
	NUMBER  // e.g. 1, 2, 3, 4
	TRUE    // true
	FALSE   // false
	NULL    // null

	COLON // :
	COMMA // ,
)

func (t Type) String() string {
	switch t {
	case LBRACKET:
		return "{"
	case RBRACKET:
		return "}"
	case LBRACE:
		return "["
	case RBRACE:
		return "]"
	case LITERAL:
		return "LITERAL"
	case NUMBER:
		return "NUMBER"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NULL:
		return "null"
	case COLON:
		return ":"
	case COMMA:
		return ","
	}
	return "unknown"
}

type Token struct {
	T Type
	V string
}

func (t Token) Type() Type {
	return t.T
}
func (t Token) Value() string {
	return t.V
}

func NewToken(t Type, v string) *Token {
	return &Token{T: t, V: v}
}

var (
	LBRACKET_TOKEN = NewToken(LBRACKET, "")
	RBRACKET_TOKEN = NewToken(RBRACKET, "")
	LBRACE_TOKEN   = NewToken(LBRACE, "")
	RBRACE_TOKEN   = NewToken(RBRACE, "")

	TRUE_TOKEN  = NewToken(TRUE, "")
	FALSE_TOKEN = NewToken(FALSE, "")
	NULL_TOKEN  = NewToken(NULL, "")

	COLON_TOKEN = NewToken(COLON, "")
	COMMA_TOKEN = NewToken(COMMA, "")
)

var KEYWORDS = map[string]*Token{
	"true":  TRUE_TOKEN,
	"false": FALSE_TOKEN,
	"null":  NULL_TOKEN,
}

func (t Token) String() string {
	return fmt.Sprintf("token(%s, %s)", t.T.String(), t.V)
}
