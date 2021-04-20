package builtins

import (
	"strings"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Split(args ...object.Object) object.Object {
	if err := typing.Check(
		"split", args,
		typing.RangeOfArgs(1, 2),
		typing.WithTypes(object.StringType, object.StringType),
	); err != nil {
		return newError(err.Error())
	}

	var sep string
	s := args[0].(*object.String).Value

	if len(args) == 2 {
		sep = args[1].(*object.String).Value
	}

	tokens := strings.Split(s, sep)
	elements := make([]object.Object, len(tokens))
	for i, token := range tokens {
		elements[i] = &object.String{Value: token}
	}
	return &object.Array{Elements: elements}
}
