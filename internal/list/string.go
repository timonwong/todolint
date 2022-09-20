package list

import (
	"strings"
)

type String []string

func NewString(items ...string) String {
	var s String
	s = append(s, items...)
	return s
}

// Set implements flag.Value interface.
func (s *String) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		*s = nil
		return nil
	}

	parts := strings.Split(v, ",")
	set := NewString(parts...)
	*s = set
	return nil
}

// String implements flag.Value interface
func (s String) String() string {
	return strings.Join(s, ",")
}
