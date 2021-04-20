package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Exit(args ...object.Object) object.Object {
	if err := typing.Check(
		"exit", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.IntegerType),
	); err != nil {
		return newError(err.Error())
	}

	var status int
	if len(args) == 1 {
		status = int(args[0].(*object.Integer).Value)
	}

	object.ExitFn(status)
	return nil
}
