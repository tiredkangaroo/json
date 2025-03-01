package lexer

import (
	"bufio"
	"io"
	"json/internals/token"
	"log"
	"strings"
	"testing"
)

const EXAMPLEJSON = `{
	"name": "John",
	"age": 27,
	"cars": [
		{
			"model_name": "Honda 2002",
			"vin": NULL,
			"years": 23,
			"needs_maintnence": true
		},
		{
			"model_name": "Water 7832",
			"vin": 107810204019401,
			"years": 5.80975e3,
			"needs_maintnence": false
		}
	]
}
`

func TestLexer(t *testing.T) {
	l := NewLexer(bufio.NewReader(strings.NewReader(EXAMPLEJSON)))
	expected := []*token.Token{
		token.LBRACKET_TOKEN, // {

		token.NewToken(token.LITERAL, "name"),
		token.COLON_TOKEN,
		token.NewToken(token.LITERAL, "John"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "age"),
		token.COLON_TOKEN,
		token.NewToken(token.NUMBER, "27"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "cars"),
		token.COLON_TOKEN,
		token.LBRACE_TOKEN, // [

		// First object in the array
		token.LBRACKET_TOKEN, // {
		token.NewToken(token.LITERAL, "model_name"),
		token.COLON_TOKEN,
		token.NewToken(token.LITERAL, "Honda 2002"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "vin"),
		token.COLON_TOKEN,
		token.NULL_TOKEN, // null
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "years"),
		token.COLON_TOKEN,
		token.NewToken(token.NUMBER, "23"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "needs_maintnence"),
		token.COLON_TOKEN,
		token.TRUE_TOKEN,     // true
		token.RBRACKET_TOKEN, // }
		token.COMMA_TOKEN,

		// Second object in the array
		token.LBRACKET_TOKEN, // {
		token.NewToken(token.LITERAL, "model_name"),
		token.COLON_TOKEN,
		token.NewToken(token.LITERAL, "Water 7832"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "vin"),
		token.COLON_TOKEN,
		token.NewToken(token.NUMBER, "107810204019401"), // undefined
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "years"),
		token.COLON_TOKEN,
		token.NewToken(token.NUMBER, "5.80975e3"),
		token.COMMA_TOKEN,

		token.NewToken(token.LITERAL, "needs_maintnence"),
		token.COLON_TOKEN,
		token.FALSE_TOKEN,    // false
		token.RBRACKET_TOKEN, // }

		token.RBRACE_TOKEN,   // ]
		token.RBRACKET_TOKEN, // }
	}

	i := 0
	for {
		tk, err := l.NextToken()
		if err == io.EOF {
			if len(expected) > i {
				t.Errorf("more tokens were expected (expected not exhausted)")
			}
			break
		}
		if err != nil {
			t.Errorf(err.Error())
			t.FailNow()
		}
		// tks = append(tks, tk)
		log.Printf("%d: %s", i, tk.String())
		if len(expected) <= i {
			t.Errorf("%d: token is not expected (expected exhausted)", i)
		} else {
			if expected[i].String() != tk.String() {
				t.Errorf("%d: not cool, expected: %s", i, expected[i].String())
			}
		}
		i++
	}
}

func BenchmarkLexer(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rd := bufio.NewReader(strings.NewReader(EXAMPLEJSON))
			l := NewLexer(rd)
			for {
				t, err := l.NextToken()
				if err == io.EOF {
					break
				}
				if err != nil {
					b.Fatalf("unexpected error: %s", err.Error())
				}
				// possibly avoids compiler-optimization
				_ = t
			}
		}
	})
	b.SetBytes(int64(len(EXAMPLEJSON)))
}
