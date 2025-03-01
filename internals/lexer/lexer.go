package lexer

import (
	"fmt"
	"io"
	"json/internals/token"
)

var TRUE_REMAINING = [...]byte{'r', 'u', 'e'}
var FALSE_REMAINING = [...]byte{'a', 'l', 's', 'e'}
var NULL_REMAINING = [...]byte{'u', 'l', 'l'}

type tokenPool struct {
	tokens      []token.Token
	lastElement int
}

func (pool *tokenPool) NewToken(t token.Type, v string) *token.Token {
	pool.lastElement++
	if pool.lastElement >= len(pool.tokens) { // it should only differ by 1
		pool.tokens = append(pool.tokens, token.Token{})
	}
	tk := &pool.tokens[pool.lastElement]
	tk.T = t
	tk.V = v
	return tk
}

type Lexer struct {
	pool *tokenPool
	*Reader
}

func (l *Lexer) NextToken() (*token.Token, error) {
	c, err := l.skipWhitespace()
	if err != nil {
		return nil, err
	}
	switch c {
	case '{':
		return token.LBRACKET_TOKEN, nil
	case '}':
		return token.RBRACKET_TOKEN, nil
	case '[':
		return token.LBRACE_TOKEN, nil
	case ']':
		return token.RBRACE_TOKEN, nil
	case ':':
		return token.COLON_TOKEN, nil
	case ',':
		return token.COMMA_TOKEN, nil
	case '"':
		return l.readLiteral()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		if isNumber(c) {
			return l.readNumber(c)
		}
	case 't':
		// possibly peek since we know how many bytes we want to read (applies to f and n)
		u := [3]byte{}
		_, err := Read(l.Reader, u[:])
		if err != nil {
			return nil, err
		}
		if u != TRUE_REMAINING {
			return nil, ErrUnknownIdentifier
		}
		return token.TRUE_TOKEN, nil
	case 'f':
		u := [4]byte{}
		_, err := Read(l.Reader, u[:])
		if err != nil {
			return nil, err
		}
		if u != FALSE_REMAINING {
			return nil, ErrUnknownIdentifier
		}
		return token.FALSE_TOKEN, nil
	case 'n':
		u := [3]byte{}
		_, err := Read(l.Reader, u[:])
		if err != nil {
			return nil, err
		}
		if u != NULL_REMAINING {
			return nil, ErrUnknownIdentifier
		}
		return token.NULL_TOKEN, nil
	}
	return nil, ErrUnknownIdentifier
}

// skipWhitespace skips the whitespace and returns the next non-whitespace character.
func (l *Lexer) skipWhitespace() (byte, error) {
	var b byte
	var err error
	for b, err = l.ReadByte(); err == nil && isWhitespace(b); b, err = l.ReadByte() {
	}
	return b, err
}

func (l *Lexer) readLiteral() (*token.Token, error) {
	s := make([]byte, 0, 12)
	var b byte
	var err error
	for b, err = l.ReadByte(); err == nil && b != '"'; b, err = l.ReadByte() {
		if b == '\\' {
			v, err := l.ReadByte() // if err != nil yes we have an unterminated literal however, throw this error for the next loop
			if err != nil {
				return nil, err
			}
			s = append(s, v) // since it's escaped just sneak it in
			continue         // we're gonna ignore the escape character
		}
		s = append(s, b)
	}
	t := l.pool.NewToken(token.LITERAL, string(s))
	return t, err
}

func (l *Lexer) readNumber(n byte) (*token.Token, error) {
	s := make([]byte, 0, 10)
	s = append(s, n)

	decimal := false
	e := false
	for {
		b, err := l.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == '-' {
			return nil, ErrInvalidNumber
		}
		if isNumber(b) {
			s = append(s, b)
			continue
		}
		if b == '.' {
			if decimal {
				return nil, ErrTooManyDecimals
			}
			if e {
				return nil, ErrInvalidScientificNotation
			}
			s = append(s, b)
			continue
		}
		if b == 'e' {
			if e {
				return nil, ErrInvalidScientificNotation
			}
			s = append(s, b)
			continue
		}
		// unread non-number byte
		l.UnreadByte()
		break
	}
	return l.pool.NewToken(token.NUMBER, string(s)), nil
}

func (l *Lexer) readKeyword(c byte) (*token.Token, error) {
	b := make([]byte, 0, 6)
	b = append(b, c)

	var j byte
	var err error
	for j, err = l.ReadByte(); isCharacter(j); j, err = l.ReadByte() {
		b = append(b, j)
	}
	// unread non-character byte
	l.UnreadByte()

	if err != nil {
		return nil, err
	}

	tk, ok := token.KEYWORDS[string(b)]
	fmt.Println("179", tk)
	if !ok {
		return nil, ErrUnknownIdentifier
	}
	return tk, nil
}

func NewLexer(rd io.Reader) Lexer {
	return Lexer{
		pool: &tokenPool{
			tokens:      make([]token.Token, 0, 256),
			lastElement: -1,
		},
		Reader: &Reader{rd: rd},
	}
}
