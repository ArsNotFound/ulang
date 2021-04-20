package builtins

import (
	"math"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func pow(x, y int64) int64 {
	p := int64(1)
	for y > 0 {
		if y&1 != 0 {
			p *= x
		}
		y >>= 1
		x *= x
	}
	return p
}

func Pow(args ...object.Object) object.Object {
	if err := typing.Check(
		"pow", args,
		typing.ExactArgs(2),
	); err != nil {
		return newError(err.Error())
	}

	switch {
	case args[0].Type() == object.IntegerType && args[1].Type() == object.IntegerType:
		x := args[0].(*object.Integer)
		y := args[1].(*object.Integer)
		value := pow(x.Value, y.Value)
		return &object.Integer{Value: value}
	case args[0].Type() == object.FloatType && args[1].Type() == object.FloatType:
		x := args[0].(*object.Float)
		y := args[1].(*object.Float)
		value := math.Pow(x.Value, y.Value)
		return &object.Float{Value: value}
	case args[0].Type() == object.IntegerType && args[1].Type() == object.FloatType:
		x := args[0].(*object.Integer)
		y := args[1].(*object.Float)
		value := math.Pow(float64(x.Value), y.Value)
		return &object.Float{Value: value}
	case args[0].Type() == object.FloatType && args[1].Type() == object.IntegerType:
		x := args[0].(*object.Float)
		y := args[1].(*object.Integer)
		value := math.Pow(x.Value, float64(y.Value))
		return &object.Float{Value: value}

	default:
		return newError("pow() takes only 'float' or 'int' as parameters")
	}
}
