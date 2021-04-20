package object

import (
	"fmt"
	"strings"
)

type HashKey struct {
	Type  Type
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Len() int {
	return len(h.Pairs)
}

func (h *Hash) Bool() bool {
	return len(h.Pairs) > 0
}

func (h *Hash) Compare(other Object) int {
	obj, ok := other.(*Hash)
	if !ok {
		return -1
	}

	if h.Len() != obj.Len() {
		return -1
	}

	for _, pair := range h.Pairs {
		left := pair.Value
		hashed := left.(Hashable)
		right, ok := obj.Pairs[hashed.HashKey()]
		if !ok {
			return -1
		}

		cmp, ok := left.(Comparable)
		if !ok {
			return -1
		}
		if cmp.Compare(right.Value) != 0 {
			return cmp.Compare(right.Value)
		}
	}

	return 0
}

func (h *Hash) String() string {
	return h.Inspect()
}

func (h *Hash) Type() Type { return HashType }

func (h *Hash) Inspect() string {
	var out strings.Builder

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteRune('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteRune('}')

	return out.String()
}
