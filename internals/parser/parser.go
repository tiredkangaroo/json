package parser

import (
	"json/internals/lexer"
)

type Parser struct {
	rd   *Reader
	root *RootNode
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
			buf: lx.PoolSlice(),
			pos: -1,
			lx:  lx,
		},
	}
}
