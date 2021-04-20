package builtins

import (
	"strconv"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Int(args ...object.Object) object.Object {
	if err := typing.Check(
		"int", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	switch arg := args[0].(type) {
	case *object.Boolean:
		return &object.Integer{Value: int64(arg.Int())}
	case *object.Integer:
		return arg
	case *object.Float:
		return &object.Integer{Value: int64(arg.Value)}
	case *object.String:
		n, err := strconv.ParseInt(arg.Value, 0, 64)
		if err != nil {
			return newError("could not parse string to int: %s", err)
		}
		return &object.Integer{Value: n}
	default:
		return newError("TypeError: cannot cast %s to int", arg.Type())
	}
}
