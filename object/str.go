package object

import (
	"fmt"
	"hash/fnv"
	"unicode/utf8"
)

type String struct {
	Value string
}

func (s *String) Len() int {
	return utf8.RuneCountInString(s.Value)
}

func (s *String) Bool() bool {
	return s.Value != ""
}

func (s *String) Compare(other Object) int {
	if obj, ok := other.(*String); ok {
		switch {
		case s.Value < obj.Value:
			return -1
		case s.Value > obj.Value:
			return 1
		default:
			return 0
		}
	}

	return 1
}

func (s *String) String() string {
	return s.Value
}

func (s *String) Inspect() string {
	return fmt.Sprintf("%#v", s.Value)
}

func (s *String) Type() Type {
	return StringType
}

func (s *String) Clone() Object {
	return &String{Value: s.Value}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
