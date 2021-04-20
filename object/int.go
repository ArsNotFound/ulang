package object

import (
	"strconv"
)

type Integer struct {
	Value int64
}

func (i *Integer) Bool() bool {
	return i.Value != 0
}

func (i *Integer) Compare(other Object) int {
	switch obj := other.(type) {
	case *Integer:
		switch {
		case i.Value < obj.Value:
			return -1
		case i.Value > obj.Value:
			return 1
		default:
			return 0
		}
	case *Float:
		switch {
		case float64(i.Value) < obj.Value:
			return -1
		case float64(i.Value) > obj.Value:
			return 1
		default:
			return 0
		}
	}

	return -1
}

func (i *Integer) String() string {
	return i.Inspect()
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Type() Type {
	return IntegerType
}

func (i *Integer) Clone() Object {
	return &Integer{Value: i.Value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
