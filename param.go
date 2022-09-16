package http

// Parameter represents a key-value pair.
type Param struct {
	Key   string
	Value string
}

// newParam creates a new parameter.
func newParam(key, value string) Param {
	return Param{key, value}
}
