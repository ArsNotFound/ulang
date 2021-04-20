package object

import (
	"strings"

	"github.com/Ars2014/ulang/ast"
)

type Function struct {
	Parameters ast.IdentifierList
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Bool() bool {
	return false
}

func (f *Function) String() string {
	return f.Inspect()
}

func (f *Function) Inspect() string {
	var out strings.Builder

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

func (f *Function) Type() Type {
	return FunctionType
}

type Return struct {
	Value Object
}

func (r *Return) Bool() bool {
	return true
}

func (r *Return) String() string {
	return r.Inspect()
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

func (r *Return) Type() Type {
	return ReturnType
}
