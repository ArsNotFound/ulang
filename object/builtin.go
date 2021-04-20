package object

import "fmt"

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Bool() bool {
	return true
}

func (b *Builtin) String() string {
	return b.Inspect()
}

func (b *Builtin) Inspect() string {
	return fmt.Sprintf("<built-in function %s>", b.Name)
}

func (b *Builtin) Type() Type {
	return BuiltInType
}
