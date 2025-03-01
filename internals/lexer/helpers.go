package lexer

import "unsafe"

func isWhitespace(b byte) bool {
	return b == '\n' || b == ' ' || b == '\t'
}

func isNumber(b byte) bool {
	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return true
	}
	return false
}

func isCharacter(b byte) bool {
	// 97 122 65 90
	return (b >= 97 && b <= 122) || (b >= 65 && b <= 90)
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
