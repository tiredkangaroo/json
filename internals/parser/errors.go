package parser

import (
	"fmt"
	"json/internals/token"
	"runtime"
)

func ErrUnexpectedEOF(expected string) error {
	return fmt.Errorf("unexpected EOF (expected: %s)", expected)
}

func LexerError(err error) error {
	return fmt.Errorf("lexer error: %s", err.Error())
}

func UnexpectedToken(expecting string, got *token.Token) error {
	callers := ""
	for i := range 4 {
		_, file, line, ok := runtime.Caller(4 - i)
		if !ok {
			continue
		}
		callers += fmt.Sprintf("%s:%d\n", file, line)
	}
	fmt.Println(callers)
	return fmt.Errorf("unexpected token: expecting %s, got: %s", expecting, got)
}
