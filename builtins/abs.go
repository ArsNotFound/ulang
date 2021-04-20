package builtins

import (
	"math"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Abs(args ...object.Object) object.Object {
	if err := typing.Check(
		"abs", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0]
	if i.Type() == object.IntegerType {
		value := i.(*object.Integer).Value
		if value < 0 {
			value = value * -1
		}
		return &object.Integer{Value: value}
	} else if i.Type() == object.FloatType {
		value := i.(*object.Float).Value
		return &object.Float{Value: math.Abs(value)}
	} else {
		return newError("TypeError: abs() expected argument to be `int` or `float` got `%s`", i.Type())
	}
}
