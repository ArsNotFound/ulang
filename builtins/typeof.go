package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func TypeOf(args ...object.Object) object.Object {
	if err := typing.Check(
		"type", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: string(args[0].Type())}
}
