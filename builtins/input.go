package builtins

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/typing"
)

func Input(args ...object.Object) object.Object {
	if err := typing.Check(
		"input", args,
		typing.RangeOfArgs(0, 1),
		typing.WithTypes(object.StringType),
	); err != nil {
		return newError(err.Error())
	}

	if len(args) == 1 {
		prompt := args[0].(*object.String).Value
		_, _ = fmt.Fprintf(object.Stdout, prompt)
	}

	buffer := bufio.NewReader(object.Stdin)

	line, _, err := buffer.ReadLine()
	if err != nil && err != io.EOF {
		return newError(fmt.Sprintf("error reading input from stdin: %s", err))
	}
	return &object.String{Value: string(line)}
}
