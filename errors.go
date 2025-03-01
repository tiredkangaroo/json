package json

import "errors"

var (
	ErrCannotParseJSONIntoGivenT = errors.New("the JSON cannot be parsed into T")
)
