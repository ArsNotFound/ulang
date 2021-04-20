package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func First(args ...object.Object) object.Object {
	if err := typing.Check(
		"first", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ArrayType),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}

	return nil
}
