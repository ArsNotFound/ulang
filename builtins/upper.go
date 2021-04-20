package builtins

import (
	"strings"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Upper(args ...object.Object) object.Object {
	if err := typing.Check(
		"upper", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.StringType),
	); err != nil {
		return newError(err.Error())
	}

	return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
}
