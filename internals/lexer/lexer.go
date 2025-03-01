package lexer

import (
	"io"
	"json/internals/token"
)

// compare
var TRUE_REMAINING = [...]byte{'r', 'u', 'e'}
var FALSE_REMAINING = [...]byte{'a', 'l', 's', 'e'}
var NULL_REMAINING = [...]byte{'u', 'l', 'l'}

type Lexer struct {
	tokens *[]token.Token
	*Reader
}

func (l *Lexer) PoolSlice() *[]token.Token {
	return l.tokens
}

func (l *Lexer) AddNewToken(t token.Type, v string) {
	*l.tokens = append(*l.tokens, token.Token{
		T: t,
		V: v,
	})
}

func (l *Lexer) AddExistingToken(t token.Token) {
	*l.tokens = append(*l.tokens, t)
}

func (l *Lexer) NextToken() error {
	c, err := l.skipWhitespace()
	if err != nil {
		return err
	}
	switch c {
	case '{':
		l.AddExistingToken(token.LBRACKET_TOKEN)
		return nil
	case '}':
		l.AddExistingToken(token.RBRACKET_TOKEN)
		return nil
	case '[':
		l.AddExistingToken(token.LBRACE_TOKEN)
		return nil
	case ']':
		l.AddExistingToken(token.RBRACE_TOKEN)
		return nil
	case ':':
		l.AddExistingToken(token.COLON_TOKEN)
		return nil
	case ',':
		l.AddExistingToken(token.COMMA_TOKEN)
		return nil
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
			return err
		}
		if u != TRUE_REMAINING {
			return ErrUnknownIdentifier
		}
		l.AddExistingToken(token.TRUE_TOKEN)
		return nil
	case 'f':
		u := [4]byte{}
		_, err := Read(l.Reader, u[:])
		if err != nil {
			return err
		}
		if u != FALSE_REMAINING {
			return ErrUnknownIdentifier
		}
		l.AddExistingToken(token.FALSE_TOKEN)
		return nil
	case 'n':
		u := [3]byte{}
		_, err := Read(l.Reader, u[:])
		if err != nil {
			return err
		}
		if u != NULL_REMAINING {
			return ErrUnknownIdentifier
		}
		l.AddExistingToken(token.NULL_TOKEN)
		return nil
	}
	return ErrUnknownIdentifier
}

// skipWhitespace skips the whitespace and returns the next non-whitespace character.
func (l *Lexer) skipWhitespace() (byte, error) {
	var b byte
	var err error
	for b, err = l.ReadByte(); err == nil && isWhitespace(b); b, err = l.ReadByte() {
	}
	return b, err
}

func (l *Lexer) readLiteral() error {
	s := make([]byte, 0, 12)
	var b byte
	var err error
	for b, err = l.ReadByte(); err == nil && b != '"'; b, err = l.ReadByte() {
		if b == '\\' {
			v, err := l.ReadByte() // if err != nil yes we have an unterminated literal however, throw this error for the next loop
			if err != nil {
				return err
			}
			s = append(s, v) // since it's escaped just sneak it in
			continue         // we're gonna ignore the escape character
		}
		s = append(s, b)
	}
	l.AddNewToken(token.LITERAL, string(s))
	return err
}

func (l *Lexer) readNumber(n byte) error {
	s := make([]byte, 0, 10)
	s = append(s, n)

	decimal := false
	e := false
	for {
		b, err := l.ReadByte()
		if err != nil {
			return err
		}
		if b == '-' {
			return ErrInvalidNumber
		}
		if isNumber(b) {
			s = append(s, b)
			continue
		}
		if b == '.' {
			if decimal {
				return ErrTooManyDecimals
			}
			if e {
				return ErrInvalidScientificNotation
			}
			s = append(s, b)
			continue
		}
		if b == 'e' {
			if e {
				return ErrInvalidScientificNotation
			}
			s = append(s, b)
			continue
		}
		// unread non-number byte
		l.UnreadByte()
		break
	}
	l.AddNewToken(token.NUMBER, string(s))
	return nil
}

// func (l *Lexer) readKeyword(c byte) (*token.Token, error) {
// 	b := make([]byte, 0, 6)
// 	b = append(b, c)

// 	var j byte
// 	var err error
// 	for j, err = l.ReadByte(); isCharacter(j); j, err = l.ReadByte() {
// 		b = append(b, j)
// 	}
// 	// unread non-character byte
// 	l.UnreadByte()

// 	if err != nil {
// 		return nil, err
// 	}

// 	tk, ok := token.KEYWORDS[string(b)]
// 	fmt.Println("179", tk)
// 	if !ok {
// 		return nil, ErrUnknownIdentifier
// 	}
// 	return tk, nil
// }

func NewLexer(rd io.Reader) Lexer {
	tks := make([]token.Token, 0, 4096)
	return Lexer{
		tokens: &tks,
		Reader: &Reader{rd: rd},
	}
}
