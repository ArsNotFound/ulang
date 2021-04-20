package builtins

import (
	"sort"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Max(args ...object.Object) object.Object {
	if err := typing.Check(
		"max", args,
		typing.MinimumArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	if args[0].Type() == object.ArrayType && len(args) == 1 {
		array := args[0].(*object.Array)
		arrayC := array.Copy()
		if len(array.Elements) < 1 {
			return nil
		}

		sort.Sort(arrayC)

		return arrayC.Elements[arrayC.Len()-1]
	} else {
		arr := &object.Array{Elements: args}
		sort.Sort(arr)
		return arr.Elements[arr.Len()-1]
	}
}
