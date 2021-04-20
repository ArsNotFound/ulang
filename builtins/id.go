package builtins

import (
	"fmt"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func IdOf(args ...object.Object) object.Object {
	if err := typing.Check(
		"id", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	arg := args[0]

	switch obj := arg.(type) {
	case *object.Null:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Boolean:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Integer:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Float:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.String:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Array:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Hash:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Function:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	case *object.Builtin:
		return &object.String{Value: fmt.Sprintf("%p", obj)}
	}

	return nil
}
