package builtins

import (
	"fmt"
	"strconv"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Hex(args ...object.Object) object.Object {
	if err := typing.Check(
		"hex", args,
		typing.ExactArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	i := args[0]
	switch i.Type() {
	case object.IntegerType:
		return &object.String{Value: fmt.Sprintf("0x%s", strconv.FormatInt(i.(*object.Integer).Value, 16))}
	case object.FloatType:
		return &object.String{Value: strconv.FormatFloat(i.(*object.Float).Value, 'x', -1, 64)}
	default:
		return newError("TypeError: hex() expected argument #1 to be 'int' or 'float', got: %s", i.Type())
	}
}
