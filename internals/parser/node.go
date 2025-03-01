package parser

import (
	"fmt"
	"io"
	"json/internals/token"
	"unsafe"
)

// ObjectNode represents a JSON Object.
type ObjectNode struct {
	// keyNodes represents the keys in the object.
	KeyNodes []KeyNode
}

// value implements the ValueNode interface for ObjectNode.
func (ObjectNode) value() {}

func (o *ObjectNode) Parse(rd *Reader) error {
	o.KeyNodes = make([]KeyNode, 0, 12)
	// keynode
	for {
		k := KeyNode{}
		err := k.Parse(rd)
		if err != nil {
			return err
		}
		o.KeyNodes = append(o.KeyNodes, k)

		t, err := rd.Peek(1)
		if err == io.EOF {
			return ErrUnexpectedEOF("key")
		}
		if err != nil {
			return LexerError(err)
		}
		if t[0] == token.RBRACKET_TOKEN { // reached end of object
			rd.Discard(1)
			break
		}
	}
	return nil
}

func (o ObjectNode) String(tabs string) string {
	s := "\n" + tabs + "Object"
	for _, k := range o.KeyNodes {
		s += k.String(tabs + "\t")
	}
	return s
}

// KeyNode represents a JSON Object Key.
type KeyNode struct {
	// Name is the name of the key.
	Name string
	// Value represents the corresponding value to the key.
	Value ValueNode
}

func (k *KeyNode) Parse(rd *Reader) error {
	// key
	// {"hello": "world"}
	//    ^
	// 	  |
	tks, err := rd.Expect(token.LITERAL, token.COLON)
	if err != nil {
		return err
	}
	// {"hello": "world"}
	//    			^
	// 	  			|
	k.Name = tks[0].Value()

	// value
	v, err := makeValueNode(rd)
	if err != nil {
		return err
	}
	k.Value = v

	// {"hello": "world", "abc": 123}
	//    				^
	// 	  				|

	// comma (we're not going to enforce JSON comma rules!)
	tks, err = rd.Peek(1)
	if err != nil {
		return err
	}
	if tks[0] == token.COMMA_TOKEN { // they used a comma
		rd.Discard(1) // consume the comma
	}

	return nil
}

func (k KeyNode) String(tabs string) string {
	s := "\n" + tabs + fmt.Sprintf("key: %s", k.Name)
	s += "\n" + tabs + fmt.Sprintf("value: %s", k.Value.String(tabs+"\t"))
	return s
}

// ArrayNode represents an JSON array.
type ArrayNode struct {
	// values represents the values in the array.
	values []ValueNode
}

// value implements the ValueNode interface for ArrayNode.
func (ArrayNode) value() {}

func (a ArrayNode) String(tabs string) string {
	s := fmt.Sprintf("\n%sArray: ", tabs)
	for _, v := range a.values {
		s += v.String(tabs + "\t")
	}
	return s
}

func (a *ArrayNode) Parse(rd *Reader) error {
	for {
		ts, err := rd.Peek(1)
		if err != nil {
			return err
		}
		t := ts[0]

		switch t {
		case token.COMMA_TOKEN: // another value
			rd.Discard(1)
			continue
		case token.RBRACE_TOKEN: // we're done here
			rd.Discard(1)
			return nil
		default:
			v, err := makeValueNode(rd)
			if err != nil {
				return err
			}
			a.values = append(a.values, v)
		}
	}
}

// StringNode represents a JSON string.
type StringNode struct {
	Value string
}

func (s StringNode) String(tabs string) string {
	return tabs + s.Value
}

// value implements the ValueNode interface for StringNode.
func (StringNode) value() {}

// BoolNode represents a JSON boolean.
type BoolNode struct {
	Value bool
}

// value implements the ValueNode interface for BoolNode.
func (BoolNode) value() {}

func (b BoolNode) String(tabs string) string {
	if b.Value {
		return tabs + "true"
	}
	return tabs + "false"
}

// NumberNode represents a JSON number.
type NumberNode struct {
	Value string
}

// value implements the ValueNode interface for NumberNode.
func (NumberNode) value() {}

func (n NumberNode) String(tabs string) string {
	return tabs + n.Value
}

// NullNode represents a JSON null value.
type NullNode struct{}

func (NullNode) value() {}

func (NullNode) String(tabs string) string {
	return tabs + "null"
}

// ValueNode is an interface that represents nodes that can be values. ValueNode
// will not be created with pointers to its values. If ValueNode is an ObjectNode,
// ValueNode.(*ObjectNode) will fail, however ValueNode.(ObjectNode) will pass.
type ValueNode interface {
	String(tabs string) string
	value()
}

func makeValueNode(rd *Reader) (ValueNode, error) {
	t, err := rd.Read()
	if err != nil {
		return nil, err
	}

	switch t.Type() {
	case token.LBRACKET: // ObjectNode
		obj := ObjectNode{}
		err := obj.Parse(rd)
		return obj, err
	case token.LBRACE:
		arr := ArrayNode{}
		err := arr.Parse(rd)
		return arr, err
	case token.LITERAL:
		str := (*StringNode)(unsafe.Pointer(&t.V))
		return str, nil
	case token.NUMBER:
		num := (*NumberNode)(unsafe.Pointer(&t.V))
		return num, nil
	case token.TRUE:
		b := BoolNode{}
		b.Value = true
		return b, nil
	case token.FALSE:
		b := BoolNode{}
		b.Value = false
		return b, nil
	case token.NULL:
		return NullNode{}, nil
	}
	return nil, UnexpectedToken("valid JSON value", t)
}
