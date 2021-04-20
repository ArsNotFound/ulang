package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Args(args ...object.Object) object.Object {
	if err := typing.Check(
		"args", args,
		typing.ExactArgs(0),
	); err != nil {
		return newError(err.Error())
	}

	elements := make([]object.Object, len(object.Arguments))
	for i, arg := range object.Arguments {
		elements[i] = &object.String{Value: arg}
	}

	return &object.Array{Elements: elements}
}
