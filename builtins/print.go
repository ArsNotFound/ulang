package builtins

import (
	"fmt"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Print(args ...object.Object) object.Object {
	if err := typing.Check(
		"print", args,
		typing.MinimumArgs(1),
	); err != nil {
		return newError(err.Error())
	}

	_, _ = fmt.Fprintln(object.Stdout, args[0].String())

	return nil
}
