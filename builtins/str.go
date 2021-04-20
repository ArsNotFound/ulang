package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Str(args ...object.Object) object.Object {
	if err := typing.Check(
		"str", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: args[0].String()}
}
