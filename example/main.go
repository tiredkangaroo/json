package main

import (
	"bufio"
	"fmt"
	"json/internals/lexer"
	"json/internals/parser"
	"log"
	"strings"
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

func main() {
	rd := bufio.NewReader(strings.NewReader(EXAMPLEJSON))
	l := lexer.NewLexer(rd)
	p := parser.NewParser(l)
	err := p.Parse()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("im done")
}
