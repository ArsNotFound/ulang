package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Last(args ...object.Object) object.Object {
	if err := typing.Check(
		"last", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ArrayType),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}

	return nil
}
