package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Push(args ...object.Object) object.Object {
	if err := typing.Check(
		"push", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.ArrayType),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.Append(args[1])
	return newArray
}
