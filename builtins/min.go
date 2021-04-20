package builtins

import (
	"sort"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Min(args ...object.Object) object.Object {
	if err := typing.Check(
		"min", args,
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

		return arrayC.Elements[0]
	} else {
		arr := &object.Array{Elements: args}
		sort.Sort(arr)
		return arr.Elements[0]
	}
}
