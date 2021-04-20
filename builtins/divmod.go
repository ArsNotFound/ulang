package builtins

import (
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Divmod(args ...object.Object) object.Object {
	if err := typing.Check(
		"divmod", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.IntegerType, object.IntegerType),
	); err != nil {
		return newError(err.Error())
	}

	a := args[0].(*object.Integer)
	b := args[0].(*object.Integer)
	elements := make([]object.Object, 2)
	elements[0] = &object.Integer{Value: a.Value / b.Value}
	elements[1] = &object.Integer{Value: a.Value % b.Value}
	return &object.Array{Elements: elements}
}
