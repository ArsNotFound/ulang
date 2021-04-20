package builtins

import (
	"fmt"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Chr(args ...object.Object) object.Object {
	if err := typing.Check(
		"chr", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.IntegerType),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	return &object.String{Value: fmt.Sprintf("%c", rune(i.Value))}
}
