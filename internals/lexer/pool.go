package lexer

import "json/internals/token"

type tokens struct {
	tokens []token.Token
}

func (tokens *tokens) NewToken(t token.Type, v string) {
	tokens.tokens = append(tokens.tokens, token.Token{
		T: t,
		V: v,
	})
}

func (tokens *tokens) AddToken(t token.Token) {
	tokens.tokens = append(tokens.tokens, t)
}

func newTokens() *tokens {
	return &tokens{
		tokens: make([]token.Token, 0, 4096),
	}
}
