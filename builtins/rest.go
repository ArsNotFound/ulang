package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Rest(args ...object.Object) object.Object {
	if err := typing.Check(
		"rest", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ArrayType),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	newArray.PopLeft()
	return newArray
}
