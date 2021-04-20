package builtins

import (
	"sort"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Sorted(args ...object.Object) object.Object {
	if err := typing.Check(
		"sort", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.ArrayType),
	); err != nil {
		return newError(err.Error())
	}

	arr := args[0].(*object.Array)
	newArray := arr.Copy()
	sort.Sort(newArray)
	return newArray
}
