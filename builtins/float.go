package builtins

import (
	"strconv"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func ToFloat(args ...object.Object) object.Object {
	if err := typing.Check(
		"float", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		return &object.Float{Value: float64(arg.Int())}
	case *object.Integer:
		return &object.Float{Value: float64(arg.Value)}
	case *object.Float:
		return arg
	case *object.String:
		n, err := strconv.ParseFloat(arg.Value, 64)
		if err != nil {
			return newError("could not parse string to int: %s", err)
		}
		return &object.Float{Value: n}
	default:
		return newError("TypeError: cannot cast %s to float", arg.Type())
	}
}
