package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Len(args ...object.Object) object.Object {
	if err := typing.Check(
		"len", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	if size, ok := args[0].(object.Sizeable); ok {
		return &object.Integer{Value: int64(size.Len())}
	}
	return newError("TypeError: object of type '%s' has no len()", args[0].Type())
}
