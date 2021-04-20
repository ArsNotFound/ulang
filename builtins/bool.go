package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Bool(args ...object.Object) object.Object {
	if err := typing.Check(
		"bool", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.Boolean{Value: args[0].Bool()}
}
