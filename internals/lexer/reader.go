package lexer

import (
	"io"
	"unsafe"
)

type Reader struct {
	rd io.Reader

	lastByte       int16
	lastByteAsNext bool
}

func (r *Reader) ReadByte() (byte, error) {
	if r.lastByteAsNext {
		r.lastByteAsNext = false
		return byte(r.lastByte), nil
	}

	var b byte
	n, err := Read(r, unsafe.Slice(&b, 1))

	if n != 1 && err == nil {
		return 0, ErrNoRead
	}

	r.lastByte = int16(b)
	return b, err
}

// UnreadByte may only follow a call to ReadByte.
func (r *Reader) UnreadByte() error {
	if r.lastByte == -1 {
		return ErrUnknownIdentifier
	}
	r.lastByteAsNext = true
	return nil
}

// first function argument is always the reciever
func readBytes(r *Reader, b []byte) (n int, err error) {
	read := 0

	n, err = r.rd.Read(b)
	read += n

	// force reading all bytes
	for read < len(b) {
		n, err = r.rd.Read(b[n:])
		read += n
	}
	r.lastByte = -1

	return n, err
}

//go:noescape
//go:linkname Read json/internals/lexer.readBytes
func Read(r *Reader, b []byte) (n int, err error)
