package parser

import (
	"fmt"
	"json/internals/lexer"
	"json/internals/token"
)

type Parser struct {
	rd   *Reader
	root *RootNode
}

type Reader struct {
	buf []*token.Token
	pos int

	lx lexer.Lexer
}

func (p *Reader) fillbuf(n int) error {
	for range n {
		tk, err := p.lx.NextToken()
		if err != nil {
			return err
		}
		p.buf = append(p.buf, tk)
	}
	return nil
}

func (p *Reader) Read() (*token.Token, error) {
	// pos starts at -1
	p.pos++
	if p.pos >= len(p.buf) {
		// find the difference between the position and the buffer
		needToFill := p.pos - len(p.buf) + 1
		if err := p.fillbuf(needToFill); err != nil {
			return nil, err
		}
	}
	return p.buf[p.pos], nil
}

func (p *Reader) Peek(n int) ([]*token.Token, error) {
	if err := p.fillbuf((p.pos + n) - len(p.buf) + 1); err != nil {
		return nil, err
	}
	return p.buf[p.pos+1 : len(p.buf)], nil
}

// Discards advances the position of the reader by n. The discard only
// has effect on the buffer on the next Read or Peek operation. This
// does NOT work like (*bufio.Reader).Discard.
func (p *Reader) Discard(n int) {
	p.pos += n
}

// Expect peeks the buffer to verify that the expected token types match the recieved
// token types from the peek. If it does not, it does not advance the reader and returns
// and error. If it does, it will advance the reader and return the recieved tokens.
func (p *Reader) Expect(expectedTks ...token.Type) ([]*token.Token, error) {
	gotTks, err := p.Peek(len(expectedTks))
	if err != nil {
		return nil, err
	}
	for i := range len(expectedTks) {
		if expectedTks[i] != gotTks[i].Type() {
			return nil, fmt.Errorf("unexpected token at %d (expected: %s, got: %s)", i, expectedTks[i].String(), gotTks[i].String())
		}
	}
	p.Discard(len(expectedTks))
	return gotTks, nil
}

func (p *Parser) Parse() error {
	p.root = new(RootNode)
	err := p.root.Parse(p.rd)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) Root() *RootNode {
	return p.root
}

func NewParser(lx lexer.Lexer) *Parser {
	return &Parser{
		rd: &Reader{
			buf: make([]*token.Token, 0, 30),
			pos: -1,
			lx:  lx,
		},
	}
}
