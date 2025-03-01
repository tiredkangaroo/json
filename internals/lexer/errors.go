package lexer

import "errors"

var (
	ErrTooManyDecimals           = errors.New("invalid number: too many decimals")
	ErrInvalidScientificNotation = errors.New("scientific notation is invalid")
	ErrInvalidNumber             = errors.New("invalid number")
	ErrUnknownIdentifier         = errors.New("unknown identifier")
	ErrNoRead                    = errors.New("n returned from io.Reader was 0; nothing happened")
	ErrUnreadOnInvalid           = errors.New("lastByte is invalid, likely because call to UnreadByte didn't follow a successful call to ReadByte")
)
