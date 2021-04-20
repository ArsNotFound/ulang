package builtins

import (
	"fmt"
	"os"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Assert(args ...object.Object) object.Object {
	if err := typing.Check(
		"assert", args,
		typing.ExactArgs(2),
		typing.WithTypes(object.BooleanType, object.StringType),
	); err != nil {
		return newError(err.Error())
	}

	if !args[0].(*object.Boolean).Value {
		fmt.Printf("Assertion Error: %s", args[1].(*object.String).Value)
		os.Exit(1)
	}

	return nil
}
