package parser

// RootNode represents the starting node for JSON. The ValueNode interface covers all possible JSON
// values, therefore the RootNode may be any JSON type.
//
// As per RFC 7159, all JSON values are valid JSON, even standalone ones (e.g. "hello world").
type RootNode struct{ v ValueNode }

func (r *RootNode) Parse(rd *Reader) (err error) {
	// a root node is just a JSON value
	r.v, err = makeValueNode(rd)
	return
}

func (r *RootNode) Value() ValueNode {
	return r.v
}

func (r *RootNode) String() string {
	return r.v.String("")
}
