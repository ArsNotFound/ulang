package builtins

import (
	"fmt"
	"strconv"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Oct(args ...object.Object) object.Object {
	if err := typing.Check(
		"oct", args,
		typing.ExactArgs(1),
		typing.WithTypes(object.IntegerType),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0].(*object.Integer)
	return &object.String{Value: fmt.Sprintf("0%s", strconv.FormatInt(i.Value, 8))}
}
