package object

import (
	"math"
	"strconv"
)

type Float struct {
	Value float64
}

func (f *Float) Bool() bool {
	return !math.IsNaN(f.Value) && f.Value != 0.0
}

func (f *Float) Compare(other Object) int {
	switch obj := other.(type) {
	case *Float:
		switch {
		case f.Value < obj.Value:
			return -1
		case f.Value > obj.Value:
			return 1
		default:
			return 0
		}
	case *Integer:
		switch {
		case f.Value < float64(obj.Value):
			return -1
		case f.Value > float64(obj.Value):
			return 1
		default:
			return 0
		}
	}

	return -1
}

func (f *Float) String() string {
	return f.Inspect()
}

func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (f *Float) Type() Type {
	return FloatType
}

func (f *Float) Clone() Object {
	return &Float{Value: f.Value}
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: math.Float64bits(f.Value)}
}
