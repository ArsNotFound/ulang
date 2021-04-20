package builtins

import (
	"strings"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Lower(args ...object.Object) object.Object {
	if err := typing.Check(
		"lower", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.StringType),
	); err != nil {
		return newError(err.Error())
	}

	str := args[0].(*object.String)
	return &object.String{Value: strings.ToLower(str.Value)}
}
